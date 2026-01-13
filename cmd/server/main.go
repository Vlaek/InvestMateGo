package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"invest-mate/internal/assets"
	"invest-mate/internal/shared/config"
	"invest-mate/internal/users"
	"invest-mate/pkg/logger"
)

type App struct {
	Config      *config.Config
	DB          *gorm.DB
	Router      *gin.Engine
	Server      *http.Server
	AssetModule *assets.Module
	UserModule  *users.Module
}

// Создание экземпляра приложения
func NewApp() *App {
	return &App{}
}

// Инициализиация приложения
func (app *App) Initialize() error {
	if err := app.loadConfiguration(); err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	app.setupGinMode()

	if err := app.setupDatabase(); err != nil {
		logger.ErrorLog("Database disabled: %v", err)
	}

	if err := app.setupModules(); err != nil {
		return fmt.Errorf("modules setup error: %w", err)
	}

	app.setupRouter()
	app.setupServer()

	return nil
}

// Загрузка конфигурации
func (app *App) loadConfiguration() error {
	if err := godotenv.Load(); err != nil {
		log.Printf("Note: .env file not found: %v", err)
	}

	app.Config = config.LoadEnv()
	return nil
}

// Настройка режима Gin в зависимости от окружения
func (app *App) setupGinMode() {
	if app.Config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
		logger.InfoLog("Running in PRODUCTION mode")
	} else {
		gin.SetMode(gin.DebugMode)
		logger.InfoLog("Running in DEBUG mode")
	}
}

// Настройка подключения к БД
func (app *App) setupDatabase() error {
	if !app.Config.IsDatabaseEnabled() {
		return fmt.Errorf("database disabled in config")
	}

	db, err := config.InitDatabase(app.Config)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB failed: %w", err)
	}

	maxOpenConns := 25
	maxIdleConns := 25
	connMaxLifetime := 5 * time.Minute

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	app.DB = db
	logger.InfoLog("Database connection established (MaxOpen: %d, MaxIdle: %d)", maxOpenConns, maxIdleConns)
	return nil
}

// Инициализация модулей приложения
func (app *App) setupModules() error {
	if app.DB == nil {
		logger.InfoLog("Running without database - modules will be initialized in limited mode")
		return nil
	}

	var err error

	app.AssetModule, err = assets.InitModule(app.DB, app.Config)
	if err != nil {
		return fmt.Errorf("assets module initialization failed: %w", err)
	}
	logger.InfoLog("Assets module initialized")

	app.UserModule, err = users.InitModule(app.DB, app.Config)
	if err != nil {
		return fmt.Errorf("users module initialization failed: %w", err)
	}
	logger.InfoLog("Users module initialized")

	return nil
}

// Установка маршрутов
func (app *App) setupRouter() {
	app.Router = gin.New()

	app.Router.Use(gin.Recovery())
	app.Router.Use(app.requestIDMiddleware())
	app.Router.Use(app.loggerMiddleware())

	app.Router.Use(app.setupCORS())
}

// Установка CORS middleware
func (app *App) setupCORS() gin.HandlerFunc {
	corsOrigins := app.Config.GetCORSOrigins()

	return cors.New(cors.Config{
		AllowOrigins: corsOrigins,
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
			"Accept",
			"X-Requested-With",
			"X-Request-ID",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
			"X-Total-Count",
			"Content-Range",
			"X-Request-ID",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowOriginFunc: func(origin string) bool {
			return app.isOriginAllowed(origin)
		},
	})
}

// Функция-проверка, разрешен ли origin
func (app *App) isOriginAllowed(origin string) bool {
	if app.Config.Env == "development" || app.Config.Env == "debug" {
		if strings.Contains(origin, "localhost") ||
			strings.Contains(origin, "127.0.0.1") ||
			strings.Contains(origin, "0.0.0.0") {
			return true
		}
	}

	for _, allowed := range app.Config.GetCORSOrigins() {
		if allowed == "*" || allowed == origin {
			return true
		}
	}

	return false
}

// Добавление Request ID к каждому запросу
func (app *App) requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = fmt.Sprintf("%d", time.Now().UnixNano())
		}
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// Логирование HTTP запросов
func (app *App) loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		if query != "" {
			path = path + "?" + query
		}

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		status := c.Writer.Status()

		logMessage := fmt.Sprintf("[GIN] %v | %3d | %13v | %15s | %-7s %s",
			c.GetString("request_id"),
			status,
			latency,
			c.ClientIP(),
			c.Request.Method,
			path,
		)

		if status >= 500 {
			logger.ErrorLog("%s", logMessage)
		} else if status >= 400 {
			logger.InfoLog("%s [CLIENT ERROR]", logMessage)
		} else {
			logger.InfoLog("%s", logMessage)
		}
	}
}

// Установка HTTP сервера
func (app *App) setupServer() {
	app.Server = &http.Server{
		Addr:         ":" + app.Config.Port,
		Handler:      app.Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

// Регистрация маршрутов приложения
func (app *App) registerRoutes() {
	api := app.Router.Group("/api/v1")

	api.GET("/health", app.healthCheck)

	if app.AssetModule != nil {
		app.AssetModule.GetHandler().RegisterRoutes(api)
		logger.InfoLog("Assets routes registered")
	}

	if app.UserModule != nil {
		app.UserModule.GetHandler().RegisterRoutes(api)
		logger.InfoLog("Users routes registered")
	}

	app.Router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":      "Not Found",
			"message":    "The requested resource does not exist",
			"path":       c.Request.URL.Path,
			"request_id": c.GetString("request_id"),
			"timestamp":  time.Now().UTC().Format(time.RFC3339),
		})
	})
}

// Обработчик запроса работоспособности приложения
func (app *App) healthCheck(c *gin.Context) {
	response := gin.H{
		"status":      "healthy",
		"service":     "invest-mate-api",
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
		"environment": app.Config.Env,
		"version":     getAppVersion(),
		"request_id":  c.GetString("request_id"),
	}

	if app.DB != nil {
		sqlDB, err := app.DB.DB()
		if err != nil {
			response["database"] = "error"
			response["status"] = "degraded"
			c.JSON(http.StatusServiceUnavailable, response)
			return
		}

		if err := sqlDB.Ping(); err != nil {
			response["database"] = "unhealthy"
			response["status"] = "degraded"
			c.JSON(http.StatusServiceUnavailable, response)
			return
		}
		response["database"] = "healthy"
	} else {
		response["database"] = "disabled"
	}

	c.JSON(http.StatusOK, response)
}

// Получение версии приложения
func getAppVersion() string {
	if version := os.Getenv("APP_VERSION"); version != "" {
		return version
	}
	return "1.0.0"
}

// Вывод информации о запуске
func (app *App) printBanner() {
	fmt.Println("══════════════════════════════════════════════════")
	fmt.Println("           InvestMate Server")
	fmt.Println("══════════════════════════════════════════════════")
	fmt.Printf("Version:     %s\n", getAppVersion())
	fmt.Printf("Environment: %s\n", app.Config.Env)
	fmt.Printf("Server Port: %s\n", app.Config.Port)
	fmt.Printf("Database:    %v\n", app.Config.IsDatabaseEnabled())

	if origins := app.Config.GetCORSOrigins(); len(origins) > 0 {
		fmt.Printf("CORS:        %v\n", origins)
	}

	fmt.Printf("Gin Mode:    %s\n", gin.Mode())
	fmt.Println("══════════════════════════════════════════════════")
}

// Запуск приложения
func (app *App) Run() error {
	app.printBanner()

	app.registerRoutes()

	serverErr := make(chan error, 1)
	go func() {
		logger.InfoLog("Starting server on port %s", app.Config.Port)
		logger.InfoLog("API available at http://localhost:%s/api/v1", app.Config.Port)

		if app.Config.Env != "production" {
			logger.InfoLog("Health check: http://localhost:%s/api/v1/health", app.Config.Port)
		}

		if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- fmt.Errorf("server error: %w", err)
		}
	}()

	select {
	case err := <-serverErr:
		return err
	case <-app.waitForShutdownSignal():
		return app.shutdown()
	}
}

// Ожидание сигналов завершения
func (app *App) waitForShutdownSignal() <-chan os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	return quit
}

// Завершение работы приложения с graceful shutdown
func (app *App) shutdown() error {
	logger.InfoLog("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	if app.DB != nil {
		if sqlDB, err := app.DB.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				logger.ErrorLog("Failed to close database connection: %v", err)
			} else {
				logger.InfoLog("Database connection closed")
			}
		}
	}

	app.closeModules()

	logger.InfoLog("Server exited cleanly")
	return nil
}

// Закрытие ресурсов модулей
func (app *App) closeModules() {
	if app.AssetModule != nil {
		if closer, ok := interface{}(app.AssetModule).(interface{ Close() error }); ok {
			if err := closer.Close(); err != nil {
				logger.ErrorLog("Failed to close assets module: %v", err)
			} else {
				logger.InfoLog("Assets module closed")
			}
		}
	}

	if app.UserModule != nil {
		if closer, ok := interface{}(app.UserModule).(interface{ Close() error }); ok {
			if err := closer.Close(); err != nil {
				logger.ErrorLog("Failed to close users module: %v", err)
			} else {
				logger.InfoLog("Users module closed")
			}
		}
	}
}

// Точка входа в приложение
func main() {
	app := NewApp()

	if err := app.Initialize(); err != nil {
		logger.ErrorLog("Failed to initialize application: %v", err)
		os.Exit(1)
	}

	if err := app.Run(); err != nil {
		logger.ErrorLog("Application error: %v", err)
		os.Exit(1)
	}
}
