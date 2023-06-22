package main

import (
	"fmt"
	scraper "scraper/pkg"
)

func main() {
	r, err := scraper.DoAPICall(1, 250)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r)
}
