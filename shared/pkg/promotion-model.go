package shared

type ProductPromotion struct {
	ID                   string `json:"promotionId"`
	TechPromoID          string `json:"techPromoId"`
	PublicationEndDate   string `json:"publicationEndDate"`
	PublicationStartDate string `json:"publicationStartDate"`
}

type Promotion struct {
	ActiveEndDate         string          `json:"activeEndDate"`
	ActiveStartDate       string          `json:"activeStartDate"`
	CommercialPromotionID string          `json:"commercialPromotionId"`
	FolderID              string          `json:"folderId"`
	MaxTimes              int             `json:"maxTimes"`
	Personalised          bool            `json:"personalised"`
	PromotionID           string          `json:"promotionId"`
	PromotionKind         string          `json:"promotionKind"`
	PromotionType         string          `json:"promotionType"`
	PublicationEndDate    string          `json:"publicationEndDate"`
	PublicationStartDate  string          `json:"publicationStartDate"`
	LinkedProducts        []LinkedProduct `json:"linkedProducts"`
	Benefit               []Benefit       `json:"benefit"`
}

type Benefit struct {
	BenefitPercentage float32 `json:"benefitPercentage"`
	MinLimit          int     `json:"minLimit"`
	LimitUnit         string  `json:"limitUnit"`
}

type LinkedProduct struct {
	ProductID              int    `json:"productId"`
	TechnicalArticleNumber string `json:"technicalArticleNumber"`
}
