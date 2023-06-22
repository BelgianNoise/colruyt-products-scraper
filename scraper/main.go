package main

import (
	"os"
	scraper "scraper/pkg"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	scraper.ScraperAPIKey = os.Getenv("SCRAPER_API_KEY")
	scraper.Scrape()
}
