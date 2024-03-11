package scraper

import "os"

var ColruytAPIEndpoint = ""
var ColruytAPIHost = ""
var ColruytPlaceID = ""

var ColruytPromotionAPIEndpoint = ""
var ColruytPromotionAPIHost = ""

func InitVariables() {
	ColruytAPIEndpoint = os.Getenv("COLRUYT_API_ENDPOINT_PRODUCTS")
	ColruytAPIHost = os.Getenv("COLRUYT_API_ENDPOINT_HOST")
	ColruytPromotionAPIEndpoint = os.Getenv("COLRUYT_API_ENDPOINT_PROMOTIONS")
	ColruytPromotionAPIHost = os.Getenv("COLRUYT_API_ENDPOINT_HOST")

	ColruytPlaceID = os.Getenv("COLRUYT_PLACE_ID")

	if ColruytAPIEndpoint == "" ||
		ColruytPromotionAPIEndpoint == "" ||
		ColruytAPIHost == "" ||
		ColruytPromotionAPIHost == "" ||
		ColruytPlaceID == "" {
		panic("Missing environment variables")
	}

}
