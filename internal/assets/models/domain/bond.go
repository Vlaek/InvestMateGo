package domain

type Bond struct {
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
	AciValue              float64 `json:"aciValue"`
	AmortizationFlag      bool    `json:"amortizationFlag"`
	CountryOfRisk         string  `json:"countryOfRisk"`
	CouponQuantityPerYear int     `json:"couponQuantityPerYear"`
	FloatingCouponFlag    bool    `json:"floatingCouponFlag"`
	InitialNominal        float64 `json:"initialNominal"`
	IssueKind             string  `json:"issueKind"`
	IssueSize             string  `json:"issueSize"`
	IssueSizePlan         string  `json:"issueSizePlan"`
	MaturityDate          string  `json:"maturityDate"`
	Nominal               float64 `json:"nominal"`
	PerpetualFlag         bool    `json:"perpetualFlag"`
	PlacementDate         string  `json:"placementDate"`
	PlacementPrice        float64 `json:"placementPrice"`
	RiskLevel             string  `json:"riskLevel"`
	Sector                string  `json:"sector"`
	StateRegDate          string  `json:"stateRegDate"`
	SubordinatedFlag      bool    `json:"subordinatedFlag"`
	BondType              string  `json:"bondType"`
}

func (b Bond) GetTicker() string { return b.Ticker }
func (b Bond) GetFigi() string   { return b.Figi }
func (b Bond) GetIsin() string   { return b.Isin }
func (b Bond) GetName() string   { return b.Name }
