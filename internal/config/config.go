package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
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

func Get() *Config {
	if AppConfig == nil {
		return LoadEnv()
	}

	return AppConfig
}

func (c *Config) GetDBDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}

func (c *Config) GetDBURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode)
}

func (c *Config) IsDBEnabled() bool {
	return c.DBHost != "" && c.DBName != ""
}

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

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")

	if valueStr == "" {
		return defaultValue
	}

	value, err := time.ParseDuration(valueStr)

	if err != nil {
		log.Printf("Invalid duration value for %s: %v", key, err)
		return defaultValue
	}

	return value
}
