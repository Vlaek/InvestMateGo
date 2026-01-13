package assets

import (
	"gorm.io/gorm"

	"invest-mate/internal/assets/handlers"
	"invest-mate/internal/assets/migrations"
	"invest-mate/internal/assets/repository"
	"invest-mate/internal/assets/services"
	"invest-mate/internal/assets/storage"
	"invest-mate/internal/shared/config"
)

type Module struct {
	assetHandler *handlers.AssetHandler
}

// Инициализация модуля
func InitModule(db *gorm.DB, cfg *config.Config) (*Module, error) {
	assetMigrator := migrations.NewAssetsMigrator()
	if err := assetMigrator.Migrate(db); err != nil {
		return nil, err
	}

	assetRepo := repository.NewAssetRepository(db)
	tinkoffStorage := storage.NewTinkoffStorage(assetRepo)
	assetService := services.NewAssetService(assetRepo, tinkoffStorage)
	assetHandler := handlers.NewAssetHandler(assetService)

	return &Module{
		assetHandler: assetHandler,
	}, nil
}
