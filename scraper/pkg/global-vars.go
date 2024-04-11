package scraper

import "os"

var ColruytAPIEndpoint = ""
var ColruytPlaceID = ""

var ColruytPromotionAPIEndpoint = ""
var ColruytPromotionAPIHost = ""

func InitVariables() {
	ColruytAPIEndpoint = os.Getenv("COLRUYT_API_ENDPOINT_PRODUCTS")
	ColruytPromotionAPIEndpoint = os.Getenv("COLRUYT_API_ENDPOINT_PROMOTIONS")
	ColruytPromotionAPIHost = os.Getenv("COLRUYT_API_ENDPOINT_HOST")

	ColruytPlaceID = os.Getenv("COLRUYT_PLACE_ID")

	if ColruytAPIEndpoint == "" {
		panic("Missing environment variable COLRUYT_API_ENDPOINT_PRODUCTS")
	}
	if ColruytPromotionAPIEndpoint == "" {
		panic("Missing environment variable COLRUYT_API_ENDPOINT_PROMOTIONS")
	}
	if ColruytPromotionAPIHost == "" {
		panic("Missing environment variable COLRUYT_API_ENDPOINT_HOST")
	}
	if ColruytPlaceID == "" {
		panic("Missing environment variable COLRUYT_PLACE_ID")
	}

}
