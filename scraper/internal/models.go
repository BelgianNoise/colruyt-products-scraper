package internal

type APIResponse struct {
	ProductsFound     int       `json:"productsFound"`
	ProductsReturned  int       `json:"productsReturned"`
	ProductsAvailable int       `json:"productsAvailable"`
	Products          []Product `json:"products"`
}

type Product struct {
	ProductID   string `json:"productId"`
	Name        string `json:"name"`
	LongName    string `json:"LongName"`
	ShortName   string `json:"ShortName"`
	Content     string `json:"content"`
	FullImage   string `json:"fullImage"`
	SquareImage string `json:"squareImage"`
	ThumbNail   string `json:"thumbNail"`
	Price       Price  `json:"price"`

	CommercialArticleNumber       string `json:"commercialArticleNumber"`
	TechnicalArticleNumber        string `json:"technicalArticleNumber"`
	AlcoholVolume                 string `json:"AlcoholVolume"`
	CountryOfOrigin               string `json:"CountryOfOrigin"`
	FicCode                       string `json:"FicCode"`
	IsBiffe                       bool   `json:"IsBiffe"`
	IsBio                         bool   `json:"IsBio"`
	IsExclusivelySoldInLuxembourg bool   `json:"IsExclusivelySoldInLuxembourg"`
	IsNew                         bool   `json:"IsNew"`
	IsPrivateLabel                bool   `json:"IsPrivateLabel"`
	IsWeightArticle               bool   `json:"IsWeightArticle"`
	OrderUnit                     string `json:"OrderUnit"`
	RecentQuanityOfStockUnits     string `json:"RecentQuanityOfStockUnits"`
	WeightconversionFactor        string `json:"WeightconversionFactor"`
	Brand                         string `json:"brand"`
	BusinessDomain                string `json:"businessDomain"`
	InPreConditionPromo           bool   `json:"inPreConditionPromo"`
	InPromo                       bool   `json:"inPromo"`
	IsAvailable                   bool   `json:"isAvailable"`
	IsPriceAvailable              bool   `json:"isPriceAvailable"`
	SeoBrand                      string `json:"seoBrand"`
	TopCategoryId                 string `json:"topCategoryId"`
	TopCategoryName               string `json:"topCategoryName"`
	WalkRouteSequenceNumber       int    `json:"walkRouteSequenceNumber"`
}

type Price struct {
	BasicPrice           float32 `json:"basicPrice"`
	IsRedPrice           bool    `json:"isRedPrice"`
	MeasurementUnit      string  `json:"measurementUnit"`
	MeasurementUnitPrice float32 `json:"measurementUnitPrice"`
	RecommendedQuantity  string  `json:"recommendedQuantity"`
}
