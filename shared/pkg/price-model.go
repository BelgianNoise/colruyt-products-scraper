package shared

type Price struct {
	BasicPrice           float32 `json:"basicPrice"`
	IsRedPrice           bool    `json:"isRedPrice"`
	MeasurementUnit      string  `json:"measurementUnit"`
	MeasurementUnitPrice float32 `json:"measurementUnitPrice"`
	RecommendedQuantity  string  `json:"recommendedQuantity"`
}
