package assets

import (
	"gorm.io/gorm"

	"invest-mate/internal/assets/handlers"
	"invest-mate/internal/assets/repository"
	"invest-mate/internal/assets/services"
	"invest-mate/internal/assets/storage"
	"invest-mate/internal/shared/config"
)

type Module struct {
	assetRepo    repository.AssetRepository
	assetService services.AssetService
	assetHandler *handlers.AssetHandler
}

// Инициализация модуля
func InitModule(db *gorm.DB, cfg *config.Config) (*Module, error) {
	assetRepo := repository.NewAssetRepository(db)

	tinkoffStorage := storage.NewTinkoffStorage(assetRepo)

	assetService := services.NewAssetService(assetRepo, tinkoffStorage)

	assetHandler := handlers.NewAssetHandler(assetService)

	return &Module{
		assetRepo:    assetRepo,
		assetService: assetService,
		assetHandler: assetHandler,
	}, nil
}

// Получение обработчика
func (m *Module) GetHandler() *handlers.AssetHandler {
	return m.assetHandler
}

// Получение сервиса
func (m *Module) GetService() services.AssetService {
	return m.assetService
}
