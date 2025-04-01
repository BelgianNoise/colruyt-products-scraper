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

	// check if the page contains an element with class "language-select-button"
	// if it not visible within 30 seconds, take a screenshot
	// if it does, click on it
	el := page.Timeout(30 * time.Second).MustElement("button.language-select-button")
	if el == nil {
		fmt.Printf("====== language-select-button not found, taking screenshot")
		page.MustScreenshot("language-select-button-not-found.png")
		panic("language-select-button not found")
	}
	el.MustClick()

	// Extract cookies from the page
	cookies = page.MustCookies()
	for _, cookie := range cookies {
		fmt.Printf("Cookie: %s=%s\n", cookie.Name, cookie.Value)
	}
}
