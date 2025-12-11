package models

type InstrumentType string

const (
    InstrumentTypeBond     InstrumentType = "BOND"
    InstrumentTypeShare    InstrumentType = "SHARE"
    InstrumentTypeETF      InstrumentType = "ETF"
    InstrumentTypeCurrency InstrumentType = "CURRENCY"
)

type Instrument interface {
    GetFigi() string
    GetTicker() string
    GetIsin() string
}