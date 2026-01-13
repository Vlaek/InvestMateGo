package users

import (
	"log"

	"gorm.io/gorm"

	"invest-mate/internal/shared/config"
	"invest-mate/internal/users/handlers"
	"invest-mate/internal/users/migrations"
	"invest-mate/internal/users/repository"
	"invest-mate/internal/users/services"
)

type Module struct {
	userRepo     repository.UserRepository
	userService  services.UserService
	userHandler  *handlers.UserHandler
	userMigrator *migrations.UsersMigrator
}

// Инициализация модуля
func InitModule(db *gorm.DB, cfg *config.Config) (*Module, error) {
	userMigrator := migrations.NewUsersMigrator()
	if err := userMigrator.Migrate(db); err != nil {
		log.Printf("⚠️ Failed to initialize users module: %v", err)
		return nil, err
	}

	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	log.Println("✅ Users module initialized")

	return &Module{
		userRepo:     userRepo,
		userService:  userService,
		userHandler:  userHandler,
		userMigrator: userMigrator,
	}, nil
}

// Получение обработчика
func (m *Module) GetHandler() *handlers.UserHandler {
	return m.userHandler
}

// Получение сервиса
func (m *Module) GetService() services.UserService {
	return m.userService
}
