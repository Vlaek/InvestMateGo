package entity

import (
	"invest-mate/internal/shared/models"
)

type Asset struct {
	Uid            string                `gorm:"primaryKey;size:255"`
	InstrumentType models.InstrumentType `gorm:"size:50"`
}
