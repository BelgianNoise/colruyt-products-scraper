package shared

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
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

func initProxyPremium() (sockets []string, err error) {
	resp, err := http.Get("https://proxypremium.top/https-ssl-proxy-list")
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	body, bodyErr := io.ReadAll(resp.Body)
	if bodyErr != nil {
		return []string{}, bodyErr
	}

	re := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+<font.*?/font>\d+`)
	ipre := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)
	portre := regexp.MustCompile(`/font>\d+`)
	matches := re.FindAllString(string(body), -1)
	for _, match := range matches {
		socket := ipre.FindString(match)
		socket += ":"
		socket += strings.TrimPrefix(portre.FindString(match), "/font>")
		sockets = append(sockets, socket)
	}
	return sockets, nil
}

func InitSockets() {
	sslProxies, err := initSslProxies()
	if err == nil {
		Sockets = append(Sockets, sslProxies...)
	}
	spys, err := initProxyPremium()
	if err == nil {
		Sockets = append(Sockets, spys...)
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
