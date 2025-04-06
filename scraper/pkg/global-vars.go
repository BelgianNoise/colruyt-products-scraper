package scraper

import (
	"fmt"
	"os"
)

var ColruytPlaceID = ""

var ColruytAPIEndpoint = "https://apip.colruyt.be/gateway/ictmgmt.emarkecom.cgproductretrsvc.v2/v2/v2/nl/products"
var ColruytPromotionAPIEndpoint = "https://apip.colruyt.be/gateway/ictmgmt.emarkecom.promotionretrsvc.v1/v1/v1/nl/promotion"

var Headless = true
var ignoreCookies = false

func InitVariables() {
	ColruytAPIEndpointEnvVar := os.Getenv("COLRUYT_API_ENDPOINT_PRODUCTS")
	if ColruytAPIEndpointEnvVar == "" {
		fmt.Printf("Using default Colruyt API endpoint: %s\n", ColruytAPIEndpoint)
	} else {
		fmt.Printf("Using Colruyt API endpoint from environment variable: %s\n", ColruytAPIEndpointEnvVar)
		ColruytAPIEndpoint = ColruytAPIEndpointEnvVar
	}
	ColruytPromotionAPIEndpointEnvVar := os.Getenv("COLRUYT_API_ENDPOINT_PROMOTIONS")
	if ColruytPromotionAPIEndpointEnvVar == "" {
		fmt.Printf("Using default Colruyt Promotion API endpoint: %s\n", ColruytPromotionAPIEndpoint)
	} else {
		fmt.Printf("Using Colruyt Promotion API endpoint from environment variable: %s\n", ColruytPromotionAPIEndpointEnvVar)
		ColruytPromotionAPIEndpoint = ColruytPromotionAPIEndpointEnvVar
	}

	ColruytPlaceID = os.Getenv("COLRUYT_PLACE_ID")
	if ColruytPlaceID == "" {
		panic("Missing environment variable COLRUYT_PLACE_ID")
	}

	HeadlessEnvVar := os.Getenv("HEADLESS")
	if HeadlessEnvVar == "" || HeadlessEnvVar == "true" {
		fmt.Printf("Using default headless mode: %v\n", Headless)
	} else {
		if HeadlessEnvVar == "false" {
			Headless = false
		}
		fmt.Printf("Using headless mode from environment variable: %v\n", Headless)
	}

	ignoreCookiesEnvVar := os.Getenv("IGNORE_COOKIES")
	if ignoreCookiesEnvVar == "" || ignoreCookiesEnvVar == "false" {
		fmt.Printf("Using default ignore cookies: %v\n", ignoreCookies)
	} else {
		if ignoreCookiesEnvVar == "true" {
			ignoreCookies = true
		}
		fmt.Printf("Using ignore cookies from environment variable: %v\n", ignoreCookies)
	}
}
