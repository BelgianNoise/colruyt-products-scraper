package shared

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

var Sockets = []string{}

func initSslProxies() (sockets []string, err error) {
	resp, err := http.Get("https://www.sslproxies.org/")
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	body, bodyErr := io.ReadAll(resp.Body)
	if bodyErr != nil {
		return []string{}, bodyErr
	}

	re := regexp.MustCompile(`<td>\d+\.\d+\.\d+\.\d+</td><td>\d+</td>`)
	matches := re.FindAllString(string(body), -1)
	for _, match := range matches {
		socket := strings.TrimPrefix(match, "<td>")
		socket = strings.TrimSuffix(socket, "</td>")
		socket = strings.Replace(socket, "</td><td>", ":", 2)
		sockets = append(sockets, socket)
	}
	return sockets, nil
}

func InitSockets() {
	sslProxies, err := initSslProxies()
	if err == nil {
		Sockets = append(Sockets, sslProxies...)
	}
	if len(Sockets) == 0 {
		panic("No proxies found, bruh we in trouble")
	}
}

var getSocketMutex = &sync.Mutex{}

func GetSocket() string {
	getSocketMutex.Lock()
	defer getSocketMutex.Unlock()
	if len(Sockets) == 0 {
		// Retrieveing new sockets makes sure we have the latest and greatest
		InitSockets()
		fmt.Printf("Sockets initialized (%d)\n", len(Sockets))
	}
	defer func() { Sockets = Sockets[1:] }()
	return Sockets[0]
}

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
