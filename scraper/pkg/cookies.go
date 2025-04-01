package scraper

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
)

var cookies []*proto.NetworkCookie = []*proto.NetworkCookie{}

func LoadCookies() {
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
