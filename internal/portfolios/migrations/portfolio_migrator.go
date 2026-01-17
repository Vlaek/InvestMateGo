package migrations

import (
	"gorm.io/gorm"

	"invest-mate/internal/portfolios/models/entity"
)

type PortfoliosMigrator struct{}

func NewPortfoliosMigrator() *PortfoliosMigrator {
	return &PortfoliosMigrator{}
}

func (m *PortfoliosMigrator) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.Portfolio{},
		&entity.Position{},
		&entity.PortfolioHierarchy{},
	)
}
