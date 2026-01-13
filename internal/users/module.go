package users

import (
	"gorm.io/gorm"

	"invest-mate/internal/shared/config"
	"invest-mate/internal/users/handlers"
	"invest-mate/internal/users/migrations"
	"invest-mate/internal/users/repository"
	"invest-mate/internal/users/services"
)

type Module struct {
	userHandler *handlers.UserHandler
}

// Инициализация модуля
func InitModule(db *gorm.DB, cfg *config.Config) (*Module, error) {
	userMigrator := migrations.NewUsersMigrator()
	if err := userMigrator.Migrate(db); err != nil {
		return nil, err
	}

	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	return &Module{
		userHandler: userHandler,
	}, nil
}
