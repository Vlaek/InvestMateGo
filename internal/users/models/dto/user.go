package dto

type PortfolioItemDTO struct {
	Figi         string  `json:"figi"`
	Ticker       string  `json:"ticker"`
	Name         string  `json:"name"`
	Quantity     float64 `json:"quantity"`
	AveragePrice float64 `json:"averagePrice"`
}

type PortfolioResponseDTO struct {
	TotalAmount float64            `json:"totalAmount"`
	Positions   []PortfolioItemDTO `json:"positions"`
}
