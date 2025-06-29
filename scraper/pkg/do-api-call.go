package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	shared "shared/pkg"
	"sync"
	"time"

	"github.com/go-rod/rod/lib/proto"
)

var userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 OPR/107.0.0.0"
var concurrencyLimiter chan int
var waitIsReloadingCookies = make(chan bool, 1)

var totalBandwidth int64 = 0
var totalBandwidthMutex = &sync.Mutex{}

func DoAPICall(
	page int,
	size int,
	useProxy bool,
	XCGAPIKey string,
) (
	responseObject APIResponse,
	err error,
) {
	// wait for the cookies reloading to be finished
	waitIsReloadingCookies <- false
	<-waitIsReloadingCookies

	var startTimeOfCall = time.Now()

	requestUrl, urlErr := url.ParseRequestURI(ColruytAPIEndpoint)
	if urlErr != nil {
		return APIResponse{}, urlErr
	}
	queryParams := requestUrl.Query()
	queryParams.Set("clientCode", "CLP")
	queryParams.Set("page", fmt.Sprint(page))
	queryParams.Set("size", fmt.Sprint(size))
	queryParams.Set("placeId", ColruytPlaceID)
	queryParams.Set("sort", "basicprice asc")
	requestUrl.RawQuery = queryParams.Encode()

	request, requestErr := http.NewRequest("GET", requestUrl.String(), nil)
	if requestErr != nil {
		return APIResponse{}, requestErr
	}

	request.Header.Set("User-Agent", userAgent)
	request.Header.Set("X-Cg-Apikey", XCGAPIKey)
	// request.Header.Set("Host", "apip.colruyt.be")
	// request.Header.Set("Origin", "https://www.colruyt.be")
	// request.Header.Set("priority", "u=1,i")
	// request.Header.Set("referer", "https://www.colruyt.be/")
	// request.Header.Set("sec-ch-ua", `"Not A(Brand";v="8", "Chromium";v="132", "Opera GX";v="117"`)
	// request.Header.Set("sec-ch-ua-mobile", "?0")
	// request.Header.Set("sec-ch-ua-platform", `"Windows"`)
	// request.Header.Set("sec-fetch-dest", "empty")
	// request.Header.Set("sec-fetch-mode", "cors")
	// request.Header.Set("sec-fetch-site", "same-site")
	// request.Header.Set("accept", "*/*")
	// request.Header.Set("accept-encoding", "gzip, deflate, br, zstd")
	// request.Header.Set("accept-language", "en")

	// Add all cookies to the request
	for _, cookie := range cookies {
		request.AddCookie(&http.Cookie{
			Name:  cookie.Name,
			Value: cookie.Value,
		})
	}

	if page == 1 && size == 1 {
		fmt.Println("- Doing initial API call")
	}

	var response *http.Response
	var responseErr error

	if useProxy {
		response, responseErr = shared.UseProxy(request)
	} else {
		client := &http.Client{Timeout: 10 * time.Second}
		response, responseErr = client.Do(request)
	}
	if responseErr != nil {
		fmt.Println(responseErr.Error())
		return retry(page, size, useProxy, XCGAPIKey)
	}
	defer response.Body.Close()

	body, bodyErr := io.ReadAll(response.Body)
	if bodyErr != nil {
		return retry(page, size, useProxy, XCGAPIKey)
	}
	var networkBandwidth = float64(len(body))
	totalBandwidthMutex.Lock()
	totalBandwidth += int64(networkBandwidth)
	totalBandwidthMutex.Unlock()

	if response.StatusCode != 200 {
		fmt.Printf("[%d] Status code: %d\n", page, response.StatusCode)
		// reload cookies if no other process is already reloading them
		if response.StatusCode == 456 && len(waitIsReloadingCookies) == 0 {
			reloadCookies()
		}
		return retry(page, size, useProxy, XCGAPIKey)
	}

	// save all cookies from the response to the global cookies variable
	for _, cookie := range response.Cookies() {
		fmt.Printf("Set Cookie: %s=%s\n", cookie.Name, cookie.Value)
		cookies = append(cookies, &proto.NetworkCookie{
			Name:  cookie.Name,
			Value: cookie.Value,
		})
	}

	var apiResponse APIResponse
	unmarshalErr := json.Unmarshal(body, &apiResponse)
	if unmarshalErr != nil {
		return retry(page, size, useProxy, XCGAPIKey)
	}

	var elapsed = time.Since(startTimeOfCall)
	fmt.Printf("[%d] Call successful (elapsed: %d ms | network bandwidth: %.2f KB)\n", page, elapsed.Milliseconds(), networkBandwidth/1024.0)

	return apiResponse, nil
}

var mayQuit = false

func retry(
	page int,
	size int,
	useProxy bool,
	XCGAPIKey string,
) (
	responseObject APIResponse,
	err error,
) {
	if mayQuit {
		return APIResponse{Products: []shared.Product{}}, nil
	}
	return DoAPICall(page, size, useProxy, XCGAPIKey)
}

func reloadCookies() {
	waitIsReloadingCookies <- true
	LoadCookies()
	<-waitIsReloadingCookies
}

func GetAllProducts() (
	products []shared.Product,
	err error,
) {
	useProxies := false
	if os.Getenv("USE_PROXY") == "true" {
		useProxies = true
	}
	return GetAllProductsWithParams(100.0, GlobalConcurrencyLimit, 250, useProxies)
}

// Retrieve a valid X-CG-APIKey.
// Not providing this header will result in a 401.
func GetXCGAPIKey() (XCGAPIKey string, err error) {
	// if the token is in the env variables, return it
	if e := os.Getenv("X_CG_APIKEY"); e != "" {
		return e, nil
	}

	// fetch https://www.colruyt.be/content/clp/nl.model.json
	// and extract the X-CG-APIKey from the response body
	request, requestErr := http.NewRequest("GET", "https://www.colruyt.be/content/clp/nl.model.json", nil)
	if requestErr != nil {
		return "", requestErr
	}
	request.Header.Set("User-Agent", userAgent)

	response, responseErr := http.DefaultClient.Do(request)
	if responseErr != nil {
		return "", responseErr
	}
	defer response.Body.Close()

	body, bodyErr := io.ReadAll(response.Body)
	if bodyErr != nil {
		return "", bodyErr
	}

	// Look for a string: '"X-CG-APIKey: a8ylmv13-b285-4788-9e14-0f79b7ed2411"'
	// and extract the key using regex
	bodyString := string(body)
	// Compile the regex pattern
	re := regexp.MustCompile(`"X-CG-APIKey: ([a-zA-Z0-9-]+)"`)

	// Find the match
	match := re.FindStringSubmatch(bodyString)
	if len(match) < 2 {
		return "", fmt.Errorf("API key not found")
	}
	apiKey := match[1]

	return apiKey, nil
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

	fmt.Printf("Using proxies: %v\n", useProxy)
	APIKey, err := GetXCGAPIKey()
	if err != nil {
		return []shared.Product{}, err
	}
	fmt.Println("API key retrieved: " + APIKey)

	percentageRequired := percentageRequiredOutOf100 / 100.0

	shared.InitProxyVars()

	initResp, err := DoAPICall(1, 1, useProxy, APIKey)
	if err != nil {
		return []shared.Product{}, err
	}

	pages := initResp.ProductsFound/pageSize + 1

	fmt.Println("making limiters for concurrency of " + fmt.Sprint(concurrencyLimit))
	concurrencyLimiter = make(chan int, concurrencyLimit)

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
			concurrencyLimiter <- 1
			wg.Add(1)
			fmt.Printf(
				"--- Acc: %d / %d (%d%s)\n",
				len(products),
				productsRequired,
				int((float32(len(products))/float32(productsRequired))*100),
				"%",
			)
			if len(products) >= int(float64(initResp.ProductsFound)*percentageRequired) {
				<-concurrencyLimiter
				wg.Done()
				fmt.Println("==========      Got enough products, breaking (pending processes will still finish)")
				mayQuit = true
				break waitTillWeGotEnoughProducts
			}
			go func(page int) {
				defer wg.Done()
				defer func() { <-concurrencyLimiter }()
				responseObject, err := DoAPICall(page, pageSize, useProxy, APIKey)
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

	fmt.Printf("[TOTAL BANDWIDTH] %.2f MB\n", float64(totalBandwidth)/1024.0/1024.0)
	return products, nil
}
