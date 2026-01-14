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
		assets.GET("/", handlers.HandleRequest(h.assetService.GetAssets))
		assets.GET("/bonds", handlers.HandleRequest(h.assetService.GetBonds))
		assets.GET("/shares", handlers.HandleRequest(h.assetService.GetShares))
		assets.GET("/etfs", handlers.HandleRequest(h.assetService.GetEtfs))
		assets.GET("/currencies", handlers.HandleRequest(h.assetService.GetCurrencies))
	}
}
