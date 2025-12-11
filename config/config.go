package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TinkoffToken   string
	Port           string
	Env            string
	LogLevel       string
	CacheTTL       int
	MaxConnections int
}

var AppConfig *Config

func Load() *Config {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	AppConfig = &Config{
		TinkoffToken:   getEnv("TINKOFF_TOKEN", ""),
		Port:           getEnv("PORT", "8080"),
		Env:            getEnv("ENV", "development"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		CacheTTL:       getEnvAsInt("CACHE_TTL", 3600),
		MaxConnections: getEnvAsInt("MAX_CONNECTIONS", 100),
	}

	if AppConfig.TinkoffToken == "" {
		log.Fatal("TINKOFF_TOKEN is required in .env file")
	}

	return AppConfig
}

func Get() *Config {
	if AppConfig == nil {
		return Load()
	}
	return AppConfig
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Invalid integer value for %s: %v", key, err)
		return defaultValue
	}

	return value
}
