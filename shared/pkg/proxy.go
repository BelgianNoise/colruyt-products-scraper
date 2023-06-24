package shared

import (
	"net/http"
	"net/url"
	"time"
)

func UseProxy(
	req *http.Request,
) (
	resp *http.Response,
	e error,
) {
	proxy := "http://" + GetSocket()
	proxyUrl, err := url.ParseRequestURI(proxy)
	if err != nil {
		return &http.Response{}, err
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
