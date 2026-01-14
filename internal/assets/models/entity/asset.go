package entity

import (
	"invest-mate/internal/shared/models"
)

type Asset struct {
	Uid            string                `gorm:"primaryKey;size:255"`
	InstrumentType models.InstrumentType `gorm:"size:50"`
}

type AssetInstrument interface {
	GetUid() string
	GetInstrumentType() models.InstrumentType
}

// Получение типа инструмента
func (a *Asset) GetInstrumentType() models.InstrumentType {
	return a.InstrumentType
}

// Получение идентификатора инструмента
func (a *Asset) GetUid() string {
	return a.Uid
}
