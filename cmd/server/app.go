package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"invest-mate/internal/shared/config"
)

type Application struct {
	Config  *config.Config
	DB      *gorm.DB
	Router  *gin.Engine
	Server  *http.Server
	Modules map[string]interface{}
}

func NewApplication() *Application {
	return &Application{
		Modules: make(map[string]interface{}),
	}
}
