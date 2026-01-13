package migrations

import (
	"gorm.io/gorm"

	"invest-mate/internal/users/models/entity"
)

type UsersMigrator struct{}

func NewUsersMigrator() *UsersMigrator {
	return &UsersMigrator{}
}

func (m *UsersMigrator) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.User{},
	)
}
