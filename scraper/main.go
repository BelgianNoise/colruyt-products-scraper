package main

import (
	scraper "scraper/pkg"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	scraper.InitVariables()
	// scraper.LoadCookies()
	scraper.Scrape()
}
