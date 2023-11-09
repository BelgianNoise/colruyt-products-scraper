package scraper

import shared "shared/pkg"

type APIResponse struct {
	ProductsFound     int              `json:"productsFound"`
	ProductsReturned  int              `json:"productsReturned"`
	ProductsAvailable int              `json:"productsAvailable"`
	Products          []shared.Product `json:"products"`
}

type PromotionAPIResponse struct {
	Promotions []shared.Promotion `json:"promotions"`
}
