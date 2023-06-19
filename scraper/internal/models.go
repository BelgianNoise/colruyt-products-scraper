package internal

type APIResponse struct {
	ProductsFound     int       `json:"productsFound"`
	ProductsReturned  int       `json:"productsReturned"`
	ProductsAvailable int       `json:"productsAvailable"`
	Products          []Product `json:"products"`
}

type Product struct {
	ProductID              string `json:"productId"`
	TechnicalArticleNumber string `json:"technicalArticleNumber"`
}
