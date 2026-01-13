package assets

import (
	"invest-mate/internal/shared/config"

	"gorm.io/gorm"
)

type ModuleWrapper struct {
	module *Module
}

func (mw *ModuleWrapper) Name() string {
	return "assets"
}

func (mw *ModuleWrapper) Initialize(db *gorm.DB, cfg *config.Config) error {
	module, err := InitModule(db, cfg)
	if err != nil {
		return err
	}
	mw.module = module
	return nil
}

func (mw *ModuleWrapper) GetHandler() interface{} {
	if mw.module == nil {
		return nil
	}
	return mw.module.GetHandler()
}

func (mw *ModuleWrapper) Close() error {
	return nil
}
