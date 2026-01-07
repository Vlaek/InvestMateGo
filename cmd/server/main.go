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

	"invest-mate/internal/config"
	"invest-mate/internal/models/entity"
	"invest-mate/internal/repository"
	"invest-mate/internal/storage"
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

	if err := db.AutoMigrate(
		&entity.Bond{},
		&entity.Share{},
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
	fmt.Println("      Tinkoff Storage Server          ")
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

	// ---------------- STORAGE ----------------
	var store *storage.TinkoffStorage

	if repo != nil {
		store = storage.NewTinkoffStorage(repo)
		log.Println("Using storage with PostgreSQL backend")
	} else {
		store = storage.GetInstanceWithoutRepo()
		log.Println("Using in-memory storage only")
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

	r.GET("/", func(c *gin.Context) {
		dbStatus := "disabled"
		if repo != nil {
			dbStatus = "connected"
		}

		c.JSON(http.StatusOK, gin.H{
			"service":  "tinkoff-storage",
			"status":   "running",
			"database": dbStatus,
			"time":     time.Now().Format(time.RFC3339),
		})
	})

	r.GET("/health", func(c *gin.Context) {
		dbStatus := "disabled"
		dbCount := 0

		if repo != nil {
			dbStatus = "connected"
		}

		bonds, _ := store.GetBonds(c.Request.Context())

		c.JSON(http.StatusOK, gin.H{
			"status":       "healthy",
			"database":     dbStatus,
			"bonds_in_mem": len(bonds),
			"bonds_in_db":  dbCount,
		})
	})

	r.GET("/bonds", func(c *gin.Context) {
		bonds, err := store.GetBonds(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"count": len(bonds),
			"bonds": bonds,
		})
	})

	r.GET("/shares", func(c *gin.Context) {
		shares, err := store.GetShares(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"count":  len(shares),
			"shares": shares,
		})
	})

	r.GET("/currencies", func(c *gin.Context) {
		currencies, err := store.GetCurrencies(c.Request.Context())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"count":      len(currencies),
			"currencies": currencies,
		})
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("🚀 Server running on http://localhost:%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// ---------------- SHUTDOWN ----------------
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("✅ Server exited cleanly")
}
