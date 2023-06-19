package main

import (
	internal "scraper/internal"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	internal.Scrape()
}
