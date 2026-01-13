package modules

import (
	"invest-mate/internal/shared/config"

	"gorm.io/gorm"
)

type BaseModule interface {
	Name() string
	Initialize(db *gorm.DB, cfg *config.Config) error
	GetHandler() interface{}
	Close() error
}

type ModuleImpl struct {
	name string
}

func (m *ModuleImpl) Name() string {
	return m.name
}

func (m *ModuleImpl) Close() error {
	return nil
}
