package main

import (
	"context"
	scraper "scraper/pkg"
	shared "shared/pkg"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	store, err := shared.NewFileStore()
	if err != nil {
		panic(err)
	}
	scraper.InitVariables()
	scraper.LoadCookies()
	scraper.Scrape(context.Background(), store)
}
