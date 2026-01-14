package models

type InstrumentType string

const (
	InstrumentTypeBond     InstrumentType = "BOND"
	InstrumentTypeShare    InstrumentType = "SHARE"
	InstrumentTypeETF      InstrumentType = "ETF"
	InstrumentTypeCurrency InstrumentType = "CURRENCY"
)
