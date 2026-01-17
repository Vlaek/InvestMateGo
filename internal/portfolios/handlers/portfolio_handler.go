package handlers

import (
	"github.com/gin-gonic/gin"

	"invest-mate/internal/portfolios/services"
)

type PortfoliosHandler struct {
	portfoliosService services.PortfoliosService
}

// Создание нового хендлера
func NewPortfoliosHandler(portfoliosService services.PortfoliosService) *PortfoliosHandler {
	return &PortfoliosHandler{portfoliosService: portfoliosService}
}

// Регистрация маршрутов
func (h *PortfoliosHandler) RegisterRoutes(router *gin.RouterGroup) {

}
