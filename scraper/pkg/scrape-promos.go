package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	shared "shared/pkg"
)

// For some reason the promotions endpoint does not have any rate limiting.
// So we don't need to do the hokey pokey with proxies, and we can just
// run sequentially.
func GetAllPromotions(
	products []shared.Product,
) (
	promotions []shared.Promotion,
	err error,
) {

	useProxies := false
	if os.Getenv("USE_PROXY") == "true" {
		useProxies = true
	}

	fmt.Printf("Using proxies: %v\n", useProxies)
	APIKey, err := GetXCGAPIKey()
	if err != nil {
		return []shared.Promotion{}, err
	}
	fmt.Println("API key retrieved: " + APIKey)

	shared.InitProxyVars()

	var promotionIDMap = map[string]string{} // Simulate a Set
	for _, prod := range products {
		for _, promotion := range prod.Promotion {
			promotionIDMap[promotion.TechPromoID] = promotion.TechPromoID
		}
	}

	fmt.Println("Promotions to scrape: ", len(promotionIDMap))

	var promotionsScraped = map[string]shared.Promotion{}
	for id, _ := range promotionIDMap {
		if _, ok := promotionsScraped[id]; !ok {
			promo, err := GetOnePromotion(id, useProxies, APIKey)
			if err == nil {
				promotionsScraped[id] = promo
				fmt.Printf("Promotion scraped: %q\n", promo.PromotionID)
			} else {
				fmt.Println("Error getting promotion: ", err)
			}
		}
	}
	for _, promo := range promotionsScraped {
		promotions = append(promotions, promo)
	}

	fmt.Println("Promotions scraped: ", len(promotions))
	return promotions, nil
}

func getOnePromotionHelper(
	promotionID string,
	useProxy bool,
	XCGAPIKey string,
) (
	promotion shared.Promotion,
	err error,
) {
	requestUrl, urlErr := url.ParseRequestURI(ColruytPromotionAPIEndpoint)
	if urlErr != nil {
		return promotion, urlErr
	}
	queryParams := requestUrl.Query()
	queryParams.Set("clientCode", "CLP")
	queryParams.Set("placeId", ColruytPlaceID)
	queryParams.Set("promotionId", promotionID)
	requestUrl.RawQuery = queryParams.Encode()

	request, requestErr := http.NewRequest("GET", requestUrl.String(), nil)
	if requestErr != nil {
		return promotion, requestErr
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", userAgent)
	request.Header.Set("X-CG-APIKey", XCGAPIKey)

	var response *http.Response
	var responseErr error

	if useProxy {
		response, responseErr = shared.UseProxy(request)
	} else {
		response, responseErr = http.DefaultClient.Do(request)
	}
	if responseErr != nil {
		return promotion, responseErr
	}

	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return promotion, err
	}

	var apiResponse shared.Promotion
	err = json.Unmarshal(bodyBytes, &apiResponse)
	if err != nil {
		return promotion, err
	}
	if apiResponse.PromotionID == "" {
		return promotion, fmt.Errorf("no promotion found for id %q", promotionID)
	}
	return apiResponse, nil
}

// Using tryCount like this only works if we are not running in goroutines.
var tryCount = 0

func GetOnePromotion(
	promorionID string,
	useProxy bool,
	XCGAPIKey string,
) (
	promotion shared.Promotion,
	err error,
) {
	tryCount++
	if tryCount > 5 {
		// Reset tryCount for next call
		tryCount = 0
		return promotion, fmt.Errorf("[%q] Tried 5 times to get promotion", promorionID)
	}
	promotion, err = getOnePromotionHelper(promorionID, useProxy, XCGAPIKey)
	if err != nil {
		fmt.Println(err.Error())
		return GetOnePromotion(promorionID, useProxy, XCGAPIKey)
	} else {
		tryCount = 0
		return promotion, nil
	}
}
