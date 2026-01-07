package domain

type InstrumentType string

const (
	InstrumentTypeBond     InstrumentType = "BOND"
	InstrumentTypeShare    InstrumentType = "SHARE"
	InstrumentTypeETF      InstrumentType = "ETF"
	InstrumentTypeCurrency InstrumentType = "CURRENCY"
)

type Instrument interface {
	GetTicker() string
	GetFigi() string
	GetIsin() string
}
