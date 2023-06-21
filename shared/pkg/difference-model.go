package shared

type PriceDifference struct {
	LongName              string  `json:"longName"`
	PriceChange           float32 `json:"priceChange"`
	PriceChangePercentage float32 `json:"priceChangePercentage"`
	InvolvesPromotion     bool    `json:"involvesPromotion"`
	OldPrice              Price   `json:"oldPrice"`
	Price                 Price   `json:"price"`
	Product
}
