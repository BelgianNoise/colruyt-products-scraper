package main

import (
	"os"
	internal "scraper/internal"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	internal.ScraperAPIKey = os.Getenv("SCRAPER_API_KEY")
	internal.Scrape()
}
