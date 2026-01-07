package users

import (
	"gorm.io/gorm"

	"invest-mate/internal/shared/config"
	"invest-mate/internal/users/handlers"
	"invest-mate/internal/users/repository"
	"invest-mate/internal/users/services"
)

// Module представляет модуль пользователей
type Module struct {
	userRepo    repository.UserRepository
	userService services.UserService
	userHandler *handlers.UserHandler
}

// NewModule создает новый модуль пользователей
func NewModule(db *gorm.DB, cfg *config.Config) (*Module, error) {
	// Инициализация репозитория
	userRepo := repository.NewUserRepository(db)

	// Инициализация сервиса
	userService := services.NewUserService(userRepo)

	// Инициализация хендлера
	userHandler := handlers.NewUserHandler(userService)

	return &Module{
		userRepo:    userRepo,
		userService: userService,
		userHandler: userHandler,
	}, nil
}

// GetHandler возвращает HTTP хендлер
func (m *Module) GetHandler() *handlers.UserHandler {
	return m.userHandler
}

// GetService возвращает сервис пользователей
func (m *Module) GetService() services.UserService {
	return m.userService
}

// Migrate выполняет миграции БД
func (m *Module) Migrate() error {
	// Создание таблиц
	// return m.userRepo.Migrate()
	return nil
}
