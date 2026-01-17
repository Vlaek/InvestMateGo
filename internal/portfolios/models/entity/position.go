package entity

import (
	"time"
)

type Position struct {
	ID                       string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	PortfolioID              string    `gorm:"not null;index;constraint:OnDelete:CASCADE"`
	Figi                     string    `gorm:"type:text;not null"`
	IsCustomData             bool      `gorm:"not null;default:false"`
	Ticker                   string    `gorm:"type:text;not null"`
	PositionUid              string    `gorm:"type:text;uniqueIndex"`
	AveragePositionPrice     float64   `gorm:"type:double precision;default:0.0"`
	AveragePositionPriceFifo float64   `gorm:"type:double precision;default:0.0"`
	AveragePositionPricePt   float64   `gorm:"type:double precision;default:0.0"`
	IsBlocked                bool      `gorm:"not null;default:false"`
	BlockedLots              int32     `gorm:"default:0"`
	CurrentPrice             float64   `gorm:"type:double precision;default:0.0"`
	DailyYield               float64   `gorm:"type:double precision;default:0.0"`
	ExpectedYield            float64   `gorm:"type:double precision;default:0.0"`
	ExpectedYieldFifo        float64   `gorm:"type:double precision;default:0.0"`
	InstrumentUid            string    `gorm:"type:text;not null"`
	Quantity                 int32     `gorm:"not null;default:0"`
	QuantityLots             int32     `gorm:"not null;default:0"`
	VarMargin                float64   `gorm:"type:double precision;default:0.0"`
	CurrentNkd               float64   `gorm:"type:double precision;default:0.0"`
	CreatedAt                time.Time `gorm:"autoCreateTime;not null"`
	UpdatedAt                time.Time `gorm:"autoUpdateTime;not null"`
}
