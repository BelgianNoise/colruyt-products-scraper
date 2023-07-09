package shared

type Promotion struct {
	ID                   string `json:"promotionId"`
	PublicationEndDate   string `json:"publicationEndDate"`
	PublicationStartDate string `json:"publicationStartDate"`
}
