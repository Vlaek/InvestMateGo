package users

import (
	"gorm.io/gorm"

	"invest-mate/internal/shared/config"
	"invest-mate/internal/users/handlers"
	"invest-mate/internal/users/repository"
	"invest-mate/internal/users/services"
)

type Module struct {
	userRepo    repository.UserRepository
	userService services.UserService
	userHandler *handlers.UserHandler
}

// Инициализация модуля
func InitModule(db *gorm.DB, cfg *config.Config) (*Module, error) {
	userRepo := repository.NewUserRepository(db)

	userService := services.NewUserService(userRepo)

	userHandler := handlers.NewUserHandler(userService)

	return &Module{
		userRepo:    userRepo,
		userService: userService,
		userHandler: userHandler,
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
