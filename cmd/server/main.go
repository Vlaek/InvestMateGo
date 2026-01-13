package main

import (
	"invest-mate/cmd/app"
	"invest-mate/pkg/logger"
	"os"
)

// Точка входа в приложение
func main() {
	app := app.NewApp()

	// Инициализация приложения
	if err := app.Initialize(); err != nil {
		logger.ErrorLog("Failed to initialize application: %v", err)
		os.Exit(1)
	}

	// Запуск приложения
	if err := app.Run(); err != nil {
		logger.ErrorLog("Application error: %v", err)
		os.Exit(1)
	}
}
