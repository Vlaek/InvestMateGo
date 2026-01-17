package portfolios

import (
	"gorm.io/gorm"

	"invest-mate/internal/portfolios/handlers"
	"invest-mate/internal/portfolios/migrations"
	"invest-mate/internal/portfolios/repository"
	"invest-mate/internal/portfolios/services"
	"invest-mate/internal/shared/config"
)

type Module struct {
	portfoliosHandler *handlers.PortfoliosHandler
}

// Инициализация модуля
func InitModule(db *gorm.DB, cfg *config.Config) (*Module, error) {
	portfoliosMigrator := migrations.NewPortfoliosMigrator()
	if err := portfoliosMigrator.Migrate(db); err != nil {
		return nil, err
	}

	portfoliosRepo := repository.NewPortfoliosRepository(db)
	portfoliosService := services.NewPortfoliosService(portfoliosRepo)
	portfoliosHandler := handlers.NewPortfoliosHandler(portfoliosService)

	return &Module{
		portfoliosHandler: portfoliosHandler,
	}, nil
}
