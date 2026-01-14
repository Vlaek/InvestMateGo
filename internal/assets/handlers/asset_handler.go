package handlers

import (
	"context"

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
		assets.GET("/", handleWithParams(h.assetService.GetAssets, h.assetService.GetAssetByField))
		assets.GET("/bonds", handleWithParams(h.assetService.GetBonds, h.assetService.GetBondByField))
		assets.GET("/shares", handleWithParams(h.assetService.GetShares, h.assetService.GetShareByField))
		assets.GET("/etfs", handleWithParams(h.assetService.GetEtfs, h.assetService.GetEtfByField))
		assets.GET("/currencies", handleWithParams(h.assetService.GetCurrencies, h.assetService.GetCurrencyByField))
	}
}

// Обработчик запроса с параметрами
func handleWithParams[T any, P any](
	getListFunc func(ctx context.Context, page, limit int) ([]T, int64, error),
	getByFieldFunc func(ctx context.Context, paramName string, paramValue string) (*P, error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.Query("uid")
		figi := c.Query("figi")
		ticker := c.Query("ticker")

		switch {
		case uid != "":
			handlers.HandleByFieldRequest(getByFieldFunc, "uid")(c)
		case figi != "":
			handlers.HandleByFieldRequest(getByFieldFunc, "figi")(c)
		case ticker != "":
			handlers.HandleByFieldRequest(getByFieldFunc, "ticker")(c)
		default:
			handlers.HandleListRequest(getListFunc)(c)
		}
	}
}
