package models

type Share struct {
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

	DivYieldFlag  bool    `json:"divYieldFlag"`
	IpoDate       string  `json:"ipoDate"`
	IssueSize     string  `json:"issueSize"`
	IssueSizePlan string  `json:"issueSizePlan"`
	Nominal       float64 `json:"nominal"`
	Sector        string  `json:"sector"`
	ShareType     string  `json:"shareType"`
}

func (s Share) GetFigi() string   { return s.Figi }
func (s Share) GetTicker() string { return s.Ticker }
func (s Share) GetIsin() string   { return s.Isin }
func (s Share) GetName() string   { return s.Name }
