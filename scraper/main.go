package main

import (
	internal "scraper/internal"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	internal.Scrape()
}
