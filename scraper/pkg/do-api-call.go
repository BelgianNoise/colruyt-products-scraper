package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	shared "shared/pkg"
	"sync"
	"time"
)

var retriesLeft = 10

func DoAPICall(
	page int,
	size int,
) (
	responseObject APIResponse,
	err error,
) {

	requestUrl, urlErr := url.ParseRequestURI(ColruytAPIEndpoint)
	if urlErr != nil {
		return APIResponse{}, urlErr
	}
	queryParams := requestUrl.Query()
	queryParams.Set("clientCode", "CLP")
	queryParams.Set("page", fmt.Sprint(page))
	queryParams.Set("size", fmt.Sprint(size))
	queryParams.Set("placeId", ColruytPlaceID)
	queryParams.Set("sort", "new desc")
	requestUrl.RawQuery = queryParams.Encode()

	scraperRequestUrl, scraperUrlErr := url.ParseRequestURI(ScraperAPIUrl)
	if scraperUrlErr != nil {
		return APIResponse{}, scraperUrlErr
	}
	scraperQueryParams := requestUrl.Query()
	scraperQueryParams.Set("api_key", ScraperAPIKey)
	scraperQueryParams.Set("keep_headers", "true")
	scraperQueryParams.Set("url", requestUrl.String())
	// scraperQueryParams.Set("render", "true")
	// scraperQueryParams.Set("session_number", "1")
	scraperQueryParams.Set("country_code", "eu")
	scraperRequestUrl.RawQuery = scraperQueryParams.Encode()

	request, requestErr := http.NewRequest("GET", scraperRequestUrl.String(), nil)
	if requestErr != nil {
		return APIResponse{}, requestErr
	}

	request.Header.Set("Host", ColruytAPIHost)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("x-cg-apikey", APIKey)
	request.Header.Set("User-Agent", UserAgent)

	fmt.Printf("[%d] Doing API call\n", page)

	response, responseErr := http.DefaultClient.Do(request)
	if responseErr != nil {
		return APIResponse{}, responseErr
	}
	defer response.Body.Close()

	fmt.Printf("[%d] Status code: %d\n", page, response.StatusCode)

	if response.StatusCode == 456 {
		if retriesLeft == 0 {
			return APIResponse{}, fmt.Errorf("API call failed")
		} else {
			retriesLeft--
			fmt.Printf("[%d] Retrying in 5 sec...\n", page)
			time.Sleep(5000 * time.Millisecond)
			return DoAPICall(page, size)
		}
	} else if response.StatusCode == 401 {
		panic("Unauthorized, Check uw key a mattie")
	}

	body, bodyErr := io.ReadAll(response.Body)
	if bodyErr != nil {
		return APIResponse{}, bodyErr
	}

	var apiResponse APIResponse
	unmarshalErr := json.Unmarshal(body, &apiResponse)
	if unmarshalErr != nil {
		if retriesLeft == 0 {
			return APIResponse{}, unmarshalErr
		} else {
			retriesLeft--
			fmt.Printf("[%d] Issue with parsing JSON, Retrying in 5 sec...\n", page)
			time.Sleep(5000 * time.Millisecond)
			return DoAPICall(page, size)
		}
	}

	return apiResponse, nil
}

func GetAllProducts() (
	products []shared.Product,
	err error,
) {

	initResp, err := DoAPICall(1, 1)
	if err != nil {
		return []shared.Product{}, err
	}
	fmt.Printf("Should retrieve %d products\n", initResp.ProductsFound)

	pages := initResp.ProductsFound/250 + 1
	repeat := 2

	// Limit to 5 concurrent requests, limit set by ScraperAPI Free plan
	limiter := make(chan int, 5)
	defer close(limiter)
	wg := sync.WaitGroup{}
	wg.Add(pages * repeat)

	productsMutex := sync.Mutex{}
	alreadyAdded := map[string]bool{}

	// For some absolute bonkers reason the API likes to go wild and return
	// different objects for the same page, so we do the same request `repeat` times.
	// It seems as if it sometimes just doesn't care about parameters passed along.
	//
	// This doesn't mean we will always get all the products, but it
	// significantly increases the % of products we get.
	//
	// Go to the `assortiment` page and order by `new`, refresh a couple of
	// times and you'll see different results, like it somehow doesn't list some
	// products. I am proper mad about this tbh.
	//
	// I could query by category to ensure that every request would only yield
	// less then 250 products. But that would lead to way too many requests,
	// considering I am on the free plan of ScraperAPI.
	for i := 1; i <= repeat; i++ {
		for i := 1; i <= pages; i++ {
			limiter <- 1
			go func(page int) {
				defer wg.Done()
				defer func() { <-limiter }()
				responseObject, err := DoAPICall(page, 250)
				if err != nil {
					fmt.Println(err)
				}

				for _, product := range responseObject.Products {
					productsMutex.Lock()
					if !alreadyAdded[product.ProductID] {
						alreadyAdded[product.ProductID] = true
						products = append(products, product)
					}
					productsMutex.Unlock()
				}

			}(i)
		}
	}

	wg.Wait()

	return products, nil

}
