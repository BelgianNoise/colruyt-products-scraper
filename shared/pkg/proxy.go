package shared

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var privateProxyHost string
var privateProxyUsername string
var privateProxyPassword string
var usePrivateProxy bool = false
var initHasBeelCalled bool = false

func InitProxyVars() {
	if initHasBeelCalled {
		return
	}
	fmt.Println("Checking proxy settings... (will use public proxies by default)")
	initHasBeelCalled = true
	if use := os.Getenv("USE_PRIVATE_PROXY"); use != "true" {
		fmt.Println("To enable private proxies, set USE_PRIVATE_PROXY=true. Continueing with public proxies...")
		return
	}
	host := os.Getenv("PRIVATE_PROXY_HOST")
	username := os.Getenv("PRIVATE_PROXY_USERNAME")
	password := os.Getenv("PRIVATE_PROXY_PASSWORD")
	if host == "" || username == "" || password == "" {
		fmt.Printf(
			"Private proxy can not be configured (host=%q,username=%q,password=%q)\n",
			host, username, strings.Repeat("*", len(password)),
		)
		return
	}
	privateProxyHost = host
	privateProxyUsername = username
	privateProxyPassword = password
	usePrivateProxy = true
	fmt.Printf(
		"Private proxy has been configured (host=%q,username=%q,password=%q)\n",
		host, username, strings.Repeat("*", len(password)),
	)
}

func UseProxy(
	req *http.Request,
) (
	resp *http.Response,
	e error,
) {
	if !initHasBeelCalled {
		panic("You are trying to use the proxy without calling InitProxyVars() first")
	}
	var proxyUrl *url.URL
	if usePrivateProxy {
		proxyUrl = &url.URL{
			Scheme: "http",
			User:   url.UserPassword(privateProxyUsername, privateProxyPassword),
			Host:   privateProxyHost,
		}
	} else {
		proxy := "http://" + GetSocket()
		proxyUrl, err := url.ParseRequestURI(proxy)
		if err != nil {
			return &http.Response{}, err
		}
		if proxyUrl == nil {
			return &http.Response{}, fmt.Errorf("proxyUrl is nil")
		}
	}

	client := &http.Client{
		Timeout: 20 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}
	defer client.CloseIdleConnections()
	response, responseErr := client.Do(req)
	if responseErr != nil {
		return &http.Response{}, responseErr
	}
	return response, nil
}
