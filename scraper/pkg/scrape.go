package scraper

import (
	"context"
	"fmt"
	shared "shared/pkg"
)

func Scrape(ctx context.Context, store *shared.FileStore) {
	fmt.Println("==========     Scraping...     ==========")

	products, err := GetAllProducts()
	if err != nil {
		panic(err)
	}

	fmt.Println("==========     Scraping done!     ==========")
	fmt.Printf("Amount of products retrieved: %d\n", len(products))
	fmt.Println("==========     Saving to storage ...     ==========")
	if err := SaveProducts(ctx, store, products); err != nil {
		panic(err)
	}
	fmt.Println("==========     Saving to storage done!     ==========")
	fmt.Println("==========     Scraping promotion data...     ==========")
	promos, err := GetAllPromotions(products)
	if err != nil {
		panic(err)
	}
	fmt.Println("==========     Scraping promotion data done!    ==========")
	fmt.Println("==========     Saving promotion data to storage ...     ==========")
	if err := SavePromotions(ctx, store, promos); err != nil {
		panic(err)
	}
	fmt.Println("==========     Saving promotion data to storage done!     ==========")
}
