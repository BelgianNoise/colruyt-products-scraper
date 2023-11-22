package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	shared "shared/pkg"
	"sync"
)

func DoAPICall(
	page int,
	size int,
	useProxy bool,
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

	request, requestErr := http.NewRequest("GET", requestUrl.String(), nil)
	if requestErr != nil {
		return APIResponse{}, requestErr
	}

	request.Header.Set("Host", ColruytAPIHost)

	if page == 1 && size == 1 {
		fmt.Println("- Doing initial API call")
	}

	var response *http.Response
	var responseErr error

	if useProxy {
		response, responseErr = shared.UseProxy(request)
	} else {
		response, responseErr = http.DefaultClient.Do(request)
	}
	if responseErr != nil {
		return retry(page, size, useProxy)
	}
	defer response.Body.Close()

	// fmt.Printf("[%d] Status code: %d\n", page, response.StatusCode)
	if response.StatusCode != 200 {
		return retry(page, size, useProxy)
	}

	body, bodyErr := io.ReadAll(response.Body)
	if bodyErr != nil {
		return retry(page, size, useProxy)
	}

	var apiResponse APIResponse
	unmarshalErr := json.Unmarshal(body, &apiResponse)
	if unmarshalErr != nil {
		return retry(page, size, useProxy)
	}

	fmt.Printf("[%d] Call successfull\n", page)

	return apiResponse, nil
}

var mayQuit = false

func retry(
	page int,
	size int,
	useProxy bool,
) (
	responseObject APIResponse,
	err error,
) {
	if mayQuit {
		return APIResponse{Products: []shared.Product{}}, nil
	}
	return DoAPICall(page, size, useProxy)
}

func GetAllProducts() (
	products []shared.Product,
	err error,
) {
	return GetAllProductsWithParams(100.0, 20, 250, true)
}

func GetAllProductsWithParams(
	percentageRequiredOutOf100 float64,
	concurrencyLimit int,
	pageSize int,
	useProxy bool,
) (
	products []shared.Product,
	err error,
) {

	percentageRequired := percentageRequiredOutOf100 / 100.0

	shared.InitProxyVars()

	initResp, err := DoAPICall(1, 1, useProxy)
	if err != nil {
		return []shared.Product{}, err
	}

	pages := initResp.ProductsFound/pageSize + 1

	limiter := make(chan int, concurrencyLimit)
	defer close(limiter)
	wg := sync.WaitGroup{}

	productsMutex := sync.Mutex{}
	alreadyAdded := map[string]bool{}

	productsRequired := int(float64(initResp.ProductsFound) * percentageRequired)
	percentRequiredString := int(percentageRequired * 100)
	fmt.Printf("Should retrieve at least %d products out of %d (%d%s)\n", productsRequired, initResp.ProductsFound, percentRequiredString, "%")

	// For some absolute bonkers reason the API likes to go wild and return
	// different objects for the same page.
	// It seems as if it sometimes just doesn't care about parameters passed along.
	//
	// Go to the `assortiment` page and order by `new`, refresh a couple of
	// times and you'll see different results, like it somehow doesn't list some
	// products. I am proper mad about this tbh.
waitTillWeGotEnoughProducts:
	for {
		for i := 1; i <= pages; i++ {
			limiter <- 1
			wg.Add(1)
			fmt.Printf(
				"--- Acc: %d / %d (%d%s)\n",
				len(products),
				productsRequired,
				int((float32(len(products))/float32(productsRequired))*100),
				"%",
			)
			if len(products) >= int(float64(initResp.ProductsFound)*percentageRequired) {
				<-limiter
				wg.Done()
				fmt.Println("==========      Got enough products, breaking (pending processes will still finish)")
				mayQuit = true
				break waitTillWeGotEnoughProducts
			}
			go func(page int) {
				defer wg.Done()
				defer func() { <-limiter }()
				responseObject, err := DoAPICall(page, pageSize, useProxy)
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
