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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"invest-mate/internal/assets"
	"invest-mate/internal/assets/models/entity"
	"invest-mate/internal/shared/config"
	"invest-mate/internal/users"
	usersEntity "invest-mate/internal/users/models/entity"
)

// TODO: Вынести из main.go
// Инициализация БД
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

	// TODO: Исправить сильную привязанность ко всем модулям
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

	var db *gorm.DB

	if cfg.IsDBEnabled() {
		var err error
		db, err = initDB(cfg)

		if err != nil {
			log.Printf("❌ DB disabled: %v", err)
			db = nil
		} else {
			sqlDB, _ := db.DB()
			defer sqlDB.Close()
		}
	}

	var assetModule *assets.Module
	var userModule *users.Module

	if db != nil {
		var err error
		assetModule, err = assets.InitModule(db, cfg)

		if err != nil {
			log.Printf("⚠️ Failed to initialize assets module: %v", err)
		} else {
			log.Println("✅ Assets module initialized")
		}

		userModule, err = users.InitModule(db, cfg)

		if err != nil {
			log.Printf("⚠️ Failed to initialize users module: %v", err)
		} else {
			log.Println("✅ Users module initialized")
		}
	} else {
		log.Println("⚠️ Nno database")
	}

	r := gin.Default()

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
			if cfg.Env == "development" {
				if strings.Contains(origin, "localhost") ||
					strings.Contains(origin, "127.0.0.1") {
					return true
				}
			}

			for _, allowed := range corsOrigins {
				if allowed == "*" || allowed == origin {
					return true
				}
			}
			return false
		},
	}))

	api := r.Group("/api/v1")

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("📚 API Documentation:")

		if assetModule != nil {
			assetHandler := assetModule.GetHandler()
			assetHandler.RegisterRoutes(api)
		}

		if userModule != nil {
			userHandler := userModule.GetHandler()
			userHandler.RegisterRoutes(api)
		}

		log.Printf("🚀 Server running on http://localhost:%s", cfg.Port)

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
