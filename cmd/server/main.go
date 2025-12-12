package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"invest-mate/config"
	"invest-mate/internal/storage"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Note: .env file not found, using environment variables: %v", err)
	}

	cfg := config.Load()

	fmt.Println("╔══════════════════════════════════════╗")
	fmt.Println("║      Tinkoff Storage Server          ║")
	fmt.Println("╚══════════════════════════════════════╝")
	fmt.Printf("Environment: %s\n", cfg.Env)
	fmt.Printf("Port: %s\n", cfg.Port)
	fmt.Printf("Token loaded: %v\n", cfg.TinkoffToken != "")
	fmt.Println("════════════════════════════════════════")

	store := storage.GetInstance()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err := store.Initialize(ctx); err != nil {
			log.Printf("Storage init failed: %v", err)
		} else {
			log.Println("✅ Storage initialized successfully")
		}
	}()

	r := gin.Default()

	// Middleware для логирования
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		log.Printf("%s %s - %v", c.Request.Method, c.Request.URL.Path, latency)
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "tinkoff-storage",
			"version": "1.0.0",
			"status":  "running",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	r.GET("/health", func(c *gin.Context) {
		if store.IsInitialized() {
			c.JSON(http.StatusOK, gin.H{"status": "healthy", "storage": "ready"})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "initializing", "storage": "loading"})
		}
	})

	r.GET("/config", func(c *gin.Context) {
		maskedToken := "****"
		if len(cfg.TinkoffToken) > 8 {
			maskedToken = cfg.TinkoffToken[:4] + "****" + cfg.TinkoffToken[len(cfg.TinkoffToken)-4:]
		}

		c.JSON(http.StatusOK, gin.H{
			"port":         cfg.Port,
			"env":          cfg.Env,
			"token_set":    cfg.TinkoffToken != "",
			"token_masked": maskedToken,
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

	addr := ":" + cfg.Port
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("🚀 Server starting on http://localhost%s", addr)
		log.Printf("📊 Endpoints:")
		log.Printf("   http://localhost%s/          - Main page", addr)
		log.Printf("   http://localhost%s/health    - Health check", addr)
		log.Printf("   http://localhost%s/config    - Config info", addr)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("🛑 Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("✅ Server exited properly")
}
