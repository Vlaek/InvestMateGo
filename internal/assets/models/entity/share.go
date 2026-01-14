package entity

import (
	"time"

	"invest-mate/internal/shared/models"
)

type Share struct {
	Uid                   string `gorm:"primaryKey;size:255"`
	Figi                  string `gorm:"size:255"`
	Ticker                string `gorm:"size:255;not null"`
	PositionUid           string `gorm:"size:255"`
	ClassCode             string `gorm:"size:255"`
	Isin                  string `gorm:"size:255"`
	Lot                   int
	Currency              string `gorm:"size:8"`
	Klong                 float64
	Kshort                float64
	Dlong                 float64
	Dshort                float64
	DlongMin              float64
	DshortMin             float64
	Exchange              string `gorm:"size:255"`
	RealExchange          string `gorm:"size:255"`
	ShortEnabledFlag      bool
	Name                  string `gorm:"size:255"`
	CountryOfRiskName     string `gorm:"size:255"`
	TradingStatus         string `gorm:"size:255"`
	OtcFlag               bool
	BuyAvailableFlag      bool
	SellAvailableFlag     bool
	MinPriceIncrement     float64
	ApiTradeAvailableFlag bool
	AssetUid              string `gorm:"size:255"`
	ForIisFlag            bool
	ForQualInvestorFlag   bool
	WeekendFlag           bool
	BlockedTcaFlag        bool
	LiquidityFlag         bool
	First1minCandleDate   string
	First1dayCandleDate   string
	DivYieldFlag          bool
	IpoDate               string
	IssueSize             string `gorm:"size:255"`
	IssueSizePlan         string `gorm:"size:255"`
	Nominal               float64
	Sector                string `gorm:"size:255"`
	ShareType             string `gorm:"size:255"`
	CreatedAt             time.Time
	UpdatedAt             time.Time

	InstrumentType models.InstrumentType `gorm:"size:50"`
}
