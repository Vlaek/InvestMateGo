package dto

import (
	"strconv"
)

type MoneyValue struct {
	Currency string `json:"currency"`
	Nano     int32  `json:"nano"`
	Units    string `json:"units"`
}

func (mv MoneyValue) ToFloat() float64 {
	return mv.ToFloatDefault(0)
}

func (mv MoneyValue) ToFloatDefault(defaultValue float64) float64 {
	units, err := strconv.ParseFloat(mv.Units, 64)
	if err != nil {
		return defaultValue
	}
	return units + float64(mv.Nano)/1e9
}
