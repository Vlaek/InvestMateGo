package domain

import "invest-mate/internal/shared/models"

type Etf struct {
	Figi                  string  `json:"figi"`
	Ticker                string  `json:"ticker"`
	PositionUid           string  `json:"positionUid"`
	ClassCode             string  `json:"classCode"`
	Isin                  string  `json:"isin"`
	Lot                   int     `json:"lot"`
	Currency              string  `json:"currency"`
	Klong                 float64 `json:"klong"`
	Kshort                float64 `json:"kshort"`
	Dlong                 float64 `json:"dlong"`
	Dshort                float64 `json:"dshort"`
	DlongMin              float64 `json:"dlongMin"`
	DshortMin             float64 `json:"dshortMin"`
	Exchange              string  `json:"exchange"`
	RealExchange          string  `json:"realExchange"`
	ShortEnabledFlag      bool    `json:"shortEnabledFlag"`
	Name                  string  `json:"name"`
	CountryOfRiskName     string  `json:"countryOfRiskName"`
	TradingStatus         string  `json:"tradingStatus"`
	OtcFlag               bool    `json:"otcFlag"`
	BuyAvailableFlag      bool    `json:"buyAvailableFlag"`
	SellAvailableFlag     bool    `json:"sellAvailableFlag"`
	MinPriceIncrement     float64 `json:"minPriceIncrement"`
	ApiTradeAvailableFlag bool    `json:"apiTradeAvailableFlag"`
	Uid                   string  `json:"uid"`
	AssetUid              string  `json:"assetUid"`
	ForIisFlag            bool    `json:"forIisFlag"`
	ForQualInvestorFlag   bool    `json:"forQualInvestorFlag"`
	WeekendFlag           bool    `json:"weekendFlag"`
	BlockedTcaFlag        bool    `json:"blockedTcaFlag"`
	LiquidityFlag         bool    `json:"liquidityFlag"`
	First1minCandleDate   string  `json:"first1minCandleDate"`
	First1dayCandleDate   string  `json:"first1dayCandleDate"`
	FixedCommission       float64 `json:"fixedCommission"`
	FocusType             string  `json:"focusType"`
	ReleasedDate          string  `json:"releasedDate"`
	NumShares             float64 `json:"numShares"`
	Sector                string  `json:"sector"`
	RebalancingFreq       string  `json:"rebalancingFreq"`

	InstrumentType models.InstrumentType `json:"instrumentType"`
}
