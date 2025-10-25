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
	if IgnoreCookies {
		fmt.Printf("====== cookies loading skipped because of env variable\n")
		return
	}

	l, browser := launchBrowser()
	defer l.Cleanup()
	defer browser.MustClose()

	page := stealth.MustPage(browser)
	page.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: userAgent,
	})

	page.MustNavigate("https://www.colruyt.be/nl/producten")
	// built in Wait load functions dont work reliably
	fmt.Println("====== waiting 10 seconds for page to load")
	time.Sleep(10 * time.Second)

	cookies = extractCookiesFromPage(page)
}

func launchBrowser() (*launcher.Launcher, *rod.Browser) {
	var browser *rod.Browser
	var l *launcher.Launcher

	fmt.Printf("====== starting browser (headless: %v)\n", Headless)

	l = launcher.New().
		Leakless(false).
		Headless(Headless).
		Devtools(!Headless).
		NoSandbox(Headless)
	url := l.MustLaunch()

	browser = rod.New().
		ControlURL(url).
		Trace(true). // Trace shows verbose debug information for each action executed
		SlowMotion(1 * time.Second).
		MustConnect()

	fmt.Printf("====== browser started")

	return l, browser
}

func extractCookiesFromPage(page *rod.Page) []*proto.NetworkCookie {
	// Extract cookies from the page
	cookies := page.MustCookies()
	for _, cookie := range cookies {
		fmt.Printf("Cookie: %s=%s\n", cookie.Name, cookie.Value)
	}
	return cookies
}
