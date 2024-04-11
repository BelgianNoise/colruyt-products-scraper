package scraper

import "os"

var ColruytAPIEndpoint = ""
var ColruytPlaceID = ""

var ColruytPromotionAPIEndpoint = ""

func InitVariables() {
	ColruytAPIEndpoint = os.Getenv("COLRUYT_API_ENDPOINT_PRODUCTS")
	ColruytPromotionAPIEndpoint = os.Getenv("COLRUYT_API_ENDPOINT_PROMOTIONS")

	ColruytPlaceID = os.Getenv("COLRUYT_PLACE_ID")

	if ColruytAPIEndpoint == "" {
		panic("Missing environment variable COLRUYT_API_ENDPOINT_PRODUCTS")
	}
	if ColruytPromotionAPIEndpoint == "" {
		panic("Missing environment variable COLRUYT_API_ENDPOINT_PROMOTIONS")
	}
	if ColruytPlaceID == "" {
		panic("Missing environment variable COLRUYT_PLACE_ID")
	}

}
