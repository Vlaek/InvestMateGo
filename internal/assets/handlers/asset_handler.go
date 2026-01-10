package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"invest-mate/internal/assets/services"
)

type AssetHandler struct {
	assetService services.AssetService
}

// Создание нового хендлера
func NewAssetHandler(assetService services.AssetService) *AssetHandler {
	return &AssetHandler{assetService: assetService}
}

// Регистрация маршрутов
func (h *AssetHandler) RegisterRoutes(router *gin.RouterGroup) {
	assets := router.Group("/assets")
	{
		assets.GET("/bonds", getAssets(h.assetService.GetBonds))
		assets.GET("/shares", getAssets(h.assetService.GetShares))
		assets.GET("/etfs", getAssets(h.assetService.GetEtfs))
		assets.GET("/currencies", getAssets(h.assetService.GetCurrencies))
	}
}

// Обобщенный хендлер для всех типов активов
func getAssets[T any](getFunc func(ctx context.Context, page, limit int) ([]T, int64, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		pageStr := c.Query("page")
		limitStr := c.Query("limit")

		page := 1
		limit := 0

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

		data, total, err := getFunc(ctx, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := gin.H{
			"data": data,
			"meta": gin.H{
				"page":  page,
				"total": total,
			},
		}

		if limit > 0 {
			pages := (int(total) + limit - 1) / limit
			response["meta"].(gin.H)["limit"] = limit
			response["meta"].(gin.H)["pages"] = pages
		}

		c.JSON(http.StatusOK, response)
	}
}
