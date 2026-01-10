package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	TinkoffToken   string
	Port           string
	Env            string
	LogLevel       string
	CacheTTL       int
	MaxConnections int

	CORSOrigins string

	DBHost         string
	DBPort         int
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	DBMaxOpenConns int
	DBMaxIdleConns int
	DBMaxIdleTime  time.Duration
}

var AppConfig *Config

// Загрузка конфигурации из .env файла
func LoadEnv() *Config {
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

		CORSOrigins: getEnv("CORS_ORIGINS", ""),

		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnvAsInt("DB_PORT", 5432),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBName:         getEnv("DB_NAME", "investmate"),
		DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
		DBMaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 25),
		DBMaxIdleTime:  time.Duration(getEnvAsInt("DB_MAX_IDLE_TIME_SECONDS", 300)) * time.Second,
	}

	if AppConfig.TinkoffToken == "" {
		log.Fatal("TINKOFF_TOKEN is required in .env file")
	}

	if AppConfig.DBPassword == "" && AppConfig.Env == "production" {
		log.Printf("Warning: DB_PASSWORD is empty in production environment")
	}

	return AppConfig
}

// Получение конфига
func GetConfig() *Config {
	if AppConfig == nil {
		return LoadEnv()
	}

	return AppConfig
}

// Проверка работы БД
func (c *Config) IsDBEnabled() bool {
	return c.DBHost != "" && c.DBName != ""
}

// Получение конфигурации БД
func (c *Config) GetDBConfig() struct {
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
} {
	return struct {
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  time.Duration
	}{
		MaxOpenConns: c.DBMaxOpenConns,
		MaxIdleConns: c.DBMaxIdleConns,
		MaxIdleTime:  c.DBMaxIdleTime,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

// Получение переменной окружения как целое число
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

// Получение CORS
func (c *Config) GetCORSOrigins() []string {
	if c.CORSOrigins == "" {
		return []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
			"http://localhost:8080",
			"http://127.0.0.1:8080",
		}
	}

	origins := strings.Split(c.CORSOrigins, ",")

	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}

	return origins
}
