package scraper

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
)

var cookies []*proto.NetworkCookie = []*proto.NetworkCookie{}

func LoadCookies() {
	var url = "https://apix.colruyt.be/gateway/ictmgmt.emarkecom.cgplacesretrsvcv4/v4/nl/places/getDetails?placeId=605"
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		panic(fmt.Sprintf("Failed to get cookies: %s", response.Status))
	}
	// log response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response: %s\n", body)

	for _, cookie := range cookies {
		cookies = append(cookies, &proto.NetworkCookie{
			Name:  cookie.Name,
			Value: cookie.Value,
		})
		fmt.Printf("Cookie Set: %s=%s\n", cookie.Name, cookie.Value)
	}
	fmt.Printf("====== cookies loaded")
}

func LoadCookiesUsingBrowser() {
	var browser *rod.Browser
	var l *launcher.Launcher

	fmt.Printf("====== starting browser")

	l = launcher.New().
		Headless(Headless).
		Devtools(!Headless).
		NoSandbox(Headless)
	defer l.Cleanup()
	url := l.MustLaunch()

	browser = rod.New().
		ControlURL(url).
		Trace(true). // Trace shows verbose debug information for each action executed
		SlowMotion(1 * time.Second).
		MustConnect()
	defer browser.MustClose()

	if !Headless {
		launcher.Open(browser.ServeMonitor(""))
	}

	fmt.Printf("====== browser started")

	page := stealth.MustPage(browser)
	page.MustNavigate("https://colruyt.be")

	var el *rod.Element
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("====== language-select-button not found, taking screenshot: %v\n", r)
				page.MustScreenshot("language-select-button-not-found.png")
				panic(r)
			}
		}()
		el = page.Timeout(30 * time.Second).MustElement(".language-select-button")
	}()

	el.MustClick()

	// Extract cookies from the page
	cookies = page.MustCookies()
	for _, cookie := range cookies {
		fmt.Printf("Cookie: %s=%s\n", cookie.Name, cookie.Value)
	}
}
