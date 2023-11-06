package shared

type Promotion struct {
	ID                   string `json:"promotionId"`
	TechPromoID          string `json:"techPromoId"`
	PublicationEndDate   string `json:"publicationEndDate"`
	PublicationStartDate string `json:"publicationStartDate"`
}
