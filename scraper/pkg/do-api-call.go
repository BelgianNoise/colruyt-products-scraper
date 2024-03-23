package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	shared "shared/pkg"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
)

var userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 OPR/107.0.0.0"

func DoAPICall(
	page int,
	size int,
	useProxy bool,
	XCGAPIKey string,
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
	queryParams.Set("sort", "basicprice asc")
	requestUrl.RawQuery = queryParams.Encode()

	request, requestErr := http.NewRequest("GET", requestUrl.String(), nil)
	if requestErr != nil {
		return APIResponse{}, requestErr
	}

	request.Header.Set("User-Agent", userAgent)
	request.Header.Set("X-Cg-Apikey", XCGAPIKey)

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
		fmt.Println(responseErr.Error())
		return retry(page, size, useProxy, XCGAPIKey)
	}
	defer response.Body.Close()

	// fmt.Printf("[%d] Status code: %d\n", page, response.StatusCode)
	if response.StatusCode != 200 {
		return retry(page, size, useProxy, XCGAPIKey)
	}

	body, bodyErr := io.ReadAll(response.Body)
	if bodyErr != nil {
		return retry(page, size, useProxy, XCGAPIKey)
	}

	var apiResponse APIResponse
	unmarshalErr := json.Unmarshal(body, &apiResponse)
	if unmarshalErr != nil {
		return retry(page, size, useProxy, XCGAPIKey)
	}

	fmt.Printf("[%d] Call successfull\n", page)

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

func GetAllProducts() (
	products []shared.Product,
	err error,
) {
	useProxies := false
	if os.Getenv("USE_PROXY") == "true" {
		useProxies = true
	}
	return GetAllProductsWithParams(100.0, 20, 250, useProxies)
}

// Retrieve a valid X-CG-APIKey from the xframe.js script.
// Not providing this header will result in a 401.
func GetXCGAPIKey() (XCGAPIKey string, err error) {
	// if the token is in the env variables, return it
	if e := os.Getenv("X_CG_APIKEY"); e != "" {
		return e, nil
	}

	var browser *rod.Browser
	var l *launcher.Launcher
	if os.Getenv("HEADLESS") == "false" {
		println("====== starting browser")

		// Headless runs the browser on foreground, you can also use flag "-rod=show"
		// Devtools opens the tab in each new tab opened automatically
		l = launcher.New().
			Headless(false).
			Devtools(true)

		defer l.Cleanup()

		url := l.MustLaunch()

		// Trace shows verbose debug information for each action executed
		// SlowMotion is a debug related function that waits 2 seconds between
		// each action, making it easier to inspect what your code is doing.
		browser = rod.New().
			ControlURL(url).
			Trace(true).
			SlowMotion(2 * time.Second).
			MustConnect()

		// ServeMonitor plays screenshots of each tab. This feature is extremely
		// useful when debugging with headless mode.
		// You can also enable it with flag "-rod=monitor"
		launcher.Open(browser.ServeMonitor(""))

		println("====== browser started")
	} else {
		browser = rod.New().SlowMotion(2 * time.Second).MustConnect()
	}
	defer browser.MustClose()

	page := stealth.MustPage(browser)
	// page := browser.MustPage("")
	router := page.HijackRequests()
	apikey := ""

	routeHandler := func(ctx *rod.Hijack) {
		// Get the value of the X-CG-APIKey header
		a := ctx.Request.Header("X-CG-APIKey")
		if a != "" && a != "<nil>" {
			fmt.Printf("API key %q found on URL %q \n", a, ctx.Request.URL().String())
			apikey = a
		} else {
			fmt.Printf("No API key found on URL %q \n", ctx.Request.URL().String())
		}
		// Continue the request
		ctx.MustLoadResponse()
	}
	router.MustAdd("*apix.colruytgroup.com*", routeHandler)
	router.MustAdd("*apip.colruytgroup.com*", routeHandler)

	go router.Run()

	page.MustNavigate("https://colruyt.be/nl")
	// wait for page to load
	page.MustWaitLoad()
	page.MustWaitRequestIdle()

	loopsToGo := 18
	for {
		// wait for 3 seconds
		time.Sleep(5 * time.Second)
		loopsToGo--
		if loopsToGo == 0 {
			break
		}
		if apikey != "" {
			return apikey, nil
		} else {
			fmt.Println("Still waiting for API key...")
		}
	}

	return "", fmt.Errorf("no API key found after 90 seconds")
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

	return products, nil

}
