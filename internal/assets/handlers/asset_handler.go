package handlers

import (
	"github.com/gin-gonic/gin"

	"invest-mate/internal/assets/services"
	"invest-mate/pkg/handlers"
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
		assets.GET("/", handlers.HandleListRequest(h.assetService.GetAssets))
		{
			assets.GET("/:uid", handlers.HandleByFieldRequest(h.assetService.GetAssetById, "uid"))
		}
		assets.GET("/bonds", handlers.HandleListRequest(h.assetService.GetBonds))
		assets.GET("/shares", handlers.HandleListRequest(h.assetService.GetShares))
		assets.GET("/etfs", handlers.HandleListRequest(h.assetService.GetEtfs))
		assets.GET("/currencies", handlers.HandleListRequest(h.assetService.GetCurrencies))
	}
}
