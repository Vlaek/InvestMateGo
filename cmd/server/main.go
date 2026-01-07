package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"invest-mate/internal/assets/models/entity"
	"invest-mate/internal/assets/repository"
	"invest-mate/internal/assets/storage"
	"invest-mate/internal/shared/config"
	"invest-mate/internal/users"
	usersEntity "invest-mate/internal/users/models/entity"
)

func initDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("gorm open: %w", err)
	}

	// Автомиграция для всех сущностей
	if err := db.AutoMigrate(
		&entity.Bond{},
		&entity.Share{},
		&entity.Etf{},
		&entity.Currency{},
		&usersEntity.User{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	log.Println("✅ PostgreSQL connected (GORM)")

	return db, nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Note: .env file not found: %v", err)
	}

	cfg := config.LoadEnv()

	fmt.Println("══════════════════════════════════════")
	fmt.Println("      InvestMate Server               ")
	fmt.Println("══════════════════════════════════════")
	fmt.Printf(
		"Env: %s | Port: %s | DB enabled: %v\n",
		cfg.Env, cfg.Port, cfg.IsDBEnabled(),
	)

	corsOrigins := cfg.GetCORSOrigins()
	fmt.Printf("CORS allowed origins: %v\n", corsOrigins)

	var (
		db   *gorm.DB
		repo *repository.PostgresRepository
	)

	if cfg.IsDBEnabled() {
		var err error
		db, err = initDB(cfg)

		if err != nil {
			log.Printf("❌ DB disabled: %v", err)
			db = nil
		} else {
			repo, _ = repository.NewPostgresRepository(cfg)
			sqlDB, _ := db.DB()
			defer sqlDB.Close()
		}
	}

	// ---------------- ИНИЦИАЛИЗАЦИЯ МОДУЛЕЙ ----------------

	// ---------------- STORAGE (Assets) ----------------
	var store *storage.TinkoffStorage

	if repo != nil {
		store = storage.NewTinkoffStorage(repo)
		log.Println("Using storage with PostgreSQL backend")
	} else {
		store = storage.GetInstanceWithoutRepo()
		log.Println("Using in-memory storage only")
	}

	// 2. Инициализация модуля Assets
	var userModule *users.Module
	if db != nil {
		var err error
		userModule, err = users.NewModule(db, cfg)
		if err != nil {
			log.Printf("⚠️ Failed to initialize users module: %v", err)
		} else {
			log.Println("✅ Users module initialized")
		}
	} else {
		log.Println("⚠️ Users module disabled (no database)")
	}

	// ---------------- INITIALIZE DATA ----------------
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err := store.Initialize(ctx); err != nil {
			log.Printf("Storage init failed: %v", err)
		} else {
			log.Println("✅ Storage initialized")
		}
	}()

	// ---------------- HTTP ----------------
	r := gin.Default()

	// Настройка CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins: corsOrigins,
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
			"Accept",
			"X-Requested-With",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
			"X-Total-Count",
			"Content-Range",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowOriginFunc: func(origin string) bool {
			// Дополнительная проверка для разработки
			if cfg.Env == "development" {
				if strings.Contains(origin, "localhost") ||
					strings.Contains(origin, "127.0.0.1") {
					return true
				}
			}

			// Проверяем разрешенные origins
			for _, allowed := range corsOrigins {
				if allowed == "*" || allowed == origin {
					return true
				}
			}
			return false
		},
	}))

	// API v1 роутер
	api := r.Group("/api/v1")

	// Регистрация маршрутов модуля Users
	if userModule != nil {
		userHandler := userModule.GetHandler()
		userHandler.RegisterRoutes(api)
		log.Println("✅ Users routes registered")
	}

	// ---------------- СУЩЕСТВУЮЩИЕ МАРШРУТЫ (Assets) ----------------
	// Эти маршруты остаются как есть, они работают через storage

	r.GET("/", func(c *gin.Context) {
		dbStatus := "disabled"
		if db != nil {
			dbStatus = "connected"
		}

		modules := gin.H{
			"users":  userModule != nil,
			"assets": true, // Assets всегда есть через storage
		}

		c.JSON(http.StatusOK, gin.H{
			"service":  "invest-mate",
			"version":  "1.0.0",
			"status":   "running",
			"database": dbStatus,
			"modules":  modules,
			"time":     time.Now().Format(time.RFC3339),
		})
	})

	r.GET("/health", func(c *gin.Context) {
		dbStatus := "disabled"
		dbCount := 0

		if db != nil {
			dbStatus = "connected"
		}

		bonds, _ := store.GetBonds(c.Request.Context())

		health := gin.H{
			"status":       "healthy",
			"database":     dbStatus,
			"bonds_in_mem": len(bonds),
			"bonds_in_db":  dbCount,
			"timestamp":    time.Now().Format(time.RFC3339),
		}

		// Добавляем информацию о модуле Users если есть
		if userModule != nil {
			health["users_module"] = "active"
		}

		c.JSON(http.StatusOK, health)
	})

	// Assets endpoints (через storage)
	r.GET("/bonds", func(c *gin.Context) {
		bonds, err := store.GetBonds(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Добавляем пагинацию
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

		total := len(bonds)
		start := (page - 1) * limit
		end := start + limit

		if start > total {
			start = total
		}
		if end > total {
			end = total
		}

		c.JSON(http.StatusOK, gin.H{
			"count": total,
			"bonds": bonds[start:end],
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": total,
				"pages": (total + limit - 1) / limit,
			},
		})
	})

	r.GET("/shares", func(c *gin.Context) {
		shares, err := store.GetShares(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Пагинация
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

		total := len(shares)
		start := (page - 1) * limit
		end := start + limit

		if start > total {
			start = total
		}
		if end > total {
			end = total
		}

		c.JSON(http.StatusOK, gin.H{
			"count":  total,
			"shares": shares[start:end],
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": total,
				"pages": (total + limit - 1) / limit,
			},
		})
	})

	r.GET("/etfs", func(c *gin.Context) {
		etfs, err := store.GetEtfs(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Пагинация
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

		total := len(etfs)
		start := (page - 1) * limit
		end := start + limit

		if start > total {
			start = total
		}
		if end > total {
			end = total
		}

		c.JSON(http.StatusOK, gin.H{
			"count": total,
			"etfs":  etfs[start:end],
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": total,
				"pages": (total + limit - 1) / limit,
			},
		})
	})

	r.GET("/currencies", func(c *gin.Context) {
		currencies, err := store.GetCurrencies(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Пагинация
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

		total := len(currencies)
		start := (page - 1) * limit
		end := start + limit

		if start > total {
			start = total
		}
		if end > total {
			end = total
		}

		c.JSON(http.StatusOK, gin.H{
			"count":      total,
			"currencies": currencies[start:end],
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": total,
				"pages": (total + limit - 1) / limit,
			},
		})
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("🚀 Server running on http://localhost:%s", cfg.Port)
		log.Printf("📚 API Endpoints:")
		log.Printf("   • Assets API:    http://localhost:%s/{bonds,shares,etfs,currencies,search}", cfg.Port)
		if userModule != nil {
			log.Printf("   • Users API:     http://localhost:%s/api/v1/users/*", cfg.Port)
		}
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// ---------------- GRACEFUL SHUTDOWN ----------------
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited cleanly")
}
