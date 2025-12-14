package models

import "time"

type Bond struct {
	Uid                   string `gorm:"primaryKey;size:255"`
	Figi                  string `gorm:"size:255"`
	Ticker                string `gorm:"size:255;not null"`
	Name                  string `gorm:"size:255"`
	Currency              string `gorm:"size:8"`
	PositionUid           string `gorm:"size:255"`
	ClassCode             string `gorm:"size:50"`
	Isin                  string `gorm:"size:50"`
	Lot                   int
	Klong                 float64
	Kshort                float64
	Dlong                 float64
	Dshort                float64
	DlongMin              float64
	DshortMin             float64
	Exchange              string `gorm:"size:50"`
	RealExchange          string `gorm:"size:50"`
	ShortEnabledFlag      bool
	CountryOfRiskName     string `gorm:"size:255"`
	TradingStatus         string `gorm:"size:50"`
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
	First1minCandleDate   string `gorm:"size:50"`
	First1dayCandleDate   string `gorm:"size:50"`
	AciValue              float64
	AmortizationFlag      bool
	CountryOfRisk         string `gorm:"size:100"`
	CouponQuantityPerYear int
	FloatingCouponFlag    bool
	InitialNominal        float64
	IssueKind             string `gorm:"size:50"`
	IssueSize             string `gorm:"size:50"`
	IssueSizePlan         string `gorm:"size:50"`
	MaturityDate          string `gorm:"size:50"`
	Nominal               float64
	PerpetualFlag         bool
	PlacementDate         string `gorm:"size:50"`
	PlacementPrice        float64
	RiskLevel             string `gorm:"size:50"`
	Sector                string `gorm:"size:100"`
	StateRegDate          string `gorm:"size:50"`
	SubordinatedFlag      bool
	BondType              string `gorm:"size:50"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
}
