package users

import (
	"os"
	"time"

	"gorm.io/gorm"

	"invest-mate/internal/shared/config"
	"invest-mate/internal/users/handlers"
	"invest-mate/internal/users/migrations"
	"invest-mate/internal/users/repository"
	"invest-mate/internal/users/services"
	middleware "invest-mate/pkg/middlewares"
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

	jwtSecret := os.Getenv("JWT_SECRET")
	middleware.InitAuthMiddleware(
		jwtSecret,
		24*time.Hour,
		7*24*time.Hour,
	)

	return &Module{
		userHandler: userHandler,
	}, nil
}
