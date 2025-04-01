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
	if !Headless {
		fmt.Printf("====== starting browser")

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

		fmt.Printf("====== browser started")
	} else {
		browser = rod.New().SlowMotion(2 * time.Second).MustConnect()
	}
	defer browser.MustClose()

	page := stealth.MustPage(browser)
	page.MustNavigate("https://colruyt.be")

	// check if the page contains an element with class "language-select-button"
	// if it does, click on it
	if page.MustElement(".language-select-button").MustVisible() {
		page.MustElement(".language-select-button").MustClick()
	}

	// Extract cookies from the page
	cookies = page.MustCookies()
	for _, cookie := range cookies {
		fmt.Printf("Cookie: %s=%s\n", cookie.Name, cookie.Value)
	}
}
