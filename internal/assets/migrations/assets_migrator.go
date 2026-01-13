package migrations

import (
	"gorm.io/gorm"

	"invest-mate/internal/assets/models/entity"
)

type AssetsMigrator struct{}

func NewAssetsMigrator() *AssetsMigrator {
	return &AssetsMigrator{}
}

func (m *AssetsMigrator) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.Bond{},
		&entity.Share{},
		&entity.Etf{},
		&entity.Currency{},
	)
}
