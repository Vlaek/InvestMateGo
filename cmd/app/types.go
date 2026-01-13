package app

import (
	"invest-mate/internal/shared/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AppModule interface {
	Initialize(db *gorm.DB, cfg *config.Config) error
	GetHandler() interface{}
	Close() error
}

type RouterRegistrar interface {
	RegisterRoutes(router *gin.RouterGroup)
}
