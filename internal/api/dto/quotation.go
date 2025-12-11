package dto

import "strconv"

type Quotation struct {
	Nano  int32  `json:"nano"`
	Units string `json:"units"`
}

func (q Quotation) ToFloat() float64 {
	return q.ToFloatDefault(0)
}

func (q Quotation) ToFloatDefault(defaultValue float64) float64 {
	units, err := strconv.ParseFloat(q.Units, 64)
	if err != nil {
		return defaultValue
	}
	return units + float64(q.Nano)/1e9
}
