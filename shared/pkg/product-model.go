package shared

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

var ProductColumns = []string{
	"id",
	"name",
	"long_name",
	"short_name",
	"content",
	"full_image",
	"square_image",
	"thumbnail",
	"commercial_article_number",
	"technical_article_number",
	"alcohol_volume",
	"country_of_origin",
	"fic_code",
	"is_biffe",
	"is_bio",
	"is_exclusively_sold_in_luxembourg",
	"is_new",
	"is_private_label",
	"is_weight_article",
	"order_unit",
	"recent_quantity_of_stock_units",
	"weightconversion_factor",
	"brand",
	"business_domain",
	"is_available",
	"seo_brand",
	"top_category_id",
	"top_category_name",
	"walk_route_sequence_number",
}
