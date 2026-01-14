package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Обобщенный хендлер для списков
func HandleListRequest[T any](getFunc func(ctx context.Context, page, limit int) ([]T, int64, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		page, limit := parsePaginationParams(c)

		data, total, err := getFunc(ctx, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := buildResponse(data, total, page, limit)
		c.JSON(http.StatusOK, response)
	}
}

func HandleByFieldRequest[T any](getFunc func(ctx context.Context, param string) (T, error), paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		param := c.Param(paramName)

		response, err := getFunc(ctx, param)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

// Функция для разбора параметров пагинации
func parsePaginationParams(c *gin.Context) (page, limit int) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page = 1
	limit = 0

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	return page, limit
}

// Функция для построения ответа с метаданными
func buildResponse(data interface{}, total int64, page, limit int) gin.H {
	meta := gin.H{
		"page":  page,
		"total": total,
	}

	if limit > 0 {
		pages := (int(total) + limit - 1) / limit
		meta["limit"] = limit
		meta["pages"] = pages
	}

	return gin.H{
		"data": data,
		"meta": meta,
	}
}
