package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	shared "shared/pkg"
)

// For some reason the promotions endpoint does not have any rate limiting.
// So we don't need to do the hokey pokey with proxies, and we can just
// run sequentially.
func GetAllPromotions(
	products []shared.Product,
) (
	promotions []shared.Promotion,
) {
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
			promo, err := GetOnePromotion(id)
			if err == nil {
				promotionsScraped[id] = promo
			}
		}
	}
	for _, promo := range promotionsScraped {
		promotions = append(promotions, promo)
	}

	fmt.Println("Promotions scraped: ", len(promotions))
	return promotions
}

func getOnePromotionHelper(
	promorionID string,
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
	queryParams.Set("promotionIds", promorionID)
	requestUrl.RawQuery = queryParams.Encode()

	request, requestErr := http.NewRequest("GET", requestUrl.String(), nil)
	if requestErr != nil {
		return promotion, requestErr
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Host", ColruytPromotionAPIHost)

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return promotion, err
	}
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return promotion, err
	}
	var apiResponse PromotionAPIResponse
	err = json.Unmarshal(bodyBytes, &apiResponse)
	if err != nil {
		return promotion, err
	}
	if len(apiResponse.Promotions) == 0 {
		return promotion, fmt.Errorf("no promotion found for id %q", promorionID)
	}
	return apiResponse.Promotions[0], nil
}

// Using tryCount like this only works if we are not running in goroutines.
var tryCount = 0

func GetOnePromotion(
	promorionID string,
) (
	promotion shared.Promotion,
	err error,
) {
	tryCount++
	if tryCount > 5 {
		return promotion, fmt.Errorf("[%q] Tried 5 times to get promotion", promorionID)
	}
	promotion, err = getOnePromotionHelper(promorionID)
	if err != nil {
		return GetOnePromotion(promorionID)
	} else {
		tryCount = 0
		return promotion, nil
	}
}
