package middleware

import (
	"fmt"
	"time"

	"invest-mate/pkg/logger"

	"github.com/gin-gonic/gin"
)

// RequestID добавляет уникальный ID к запросу
func RequestID() gin.HandlerFunc {
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

// LoggerMiddleware логирует HTTP запросы
func LoggerMiddleware() gin.HandlerFunc {
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
		requestID := c.GetString("request_id")

		// Форматируем логирование
		logMessage := fmt.Sprintf("[GIN] %v | %3d | %13v | %15s | %-7s %s",
			requestID,
			status,
			latency,
			c.ClientIP(),
			c.Request.Method,
			path,
		)

		// Разделяем логи по статусам
		if status >= 500 {
			logger.ErrorLog(logMessage)
		} else if status >= 400 {
			logger.InfoLog(logMessage + " [CLIENT ERROR]")
		} else {
			logger.InfoLog(logMessage)
		}
	}
}
