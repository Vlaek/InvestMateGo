package assets

import (
	"log"

	"gorm.io/gorm"

	"invest-mate/internal/assets/handlers"
	"invest-mate/internal/assets/migrations"
	"invest-mate/internal/assets/repository"
	"invest-mate/internal/assets/services"
	"invest-mate/internal/assets/storage"
	"invest-mate/internal/shared/config"
)

type Module struct {
	assetRepo     repository.AssetRepository
	assetService  services.AssetService
	assetHandler  *handlers.AssetHandler
	assetMigrator *migrations.AssetsMigrator
}

// Инициализация модуля
func InitModule(db *gorm.DB, cfg *config.Config) (*Module, error) {
	assetMigrator := migrations.NewAssetsMigrator()
	if err := assetMigrator.Migrate(db); err != nil {
		log.Printf("⚠️ Failed to initialize assets module: %v", err)
		return nil, err
	}

	assetRepo := repository.NewAssetRepository(db)
	tinkoffStorage := storage.NewTinkoffStorage(assetRepo)
	assetService := services.NewAssetService(assetRepo, tinkoffStorage)
	assetHandler := handlers.NewAssetHandler(assetService)

	log.Println("✅ Assets module initialized")

	return &Module{
		assetRepo:     assetRepo,
		assetService:  assetService,
		assetHandler:  assetHandler,
		assetMigrator: assetMigrator,
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

// Получение мигратора
func (m *Module) GetMigrator() *migrations.AssetsMigrator {
	return m.assetMigrator
}
