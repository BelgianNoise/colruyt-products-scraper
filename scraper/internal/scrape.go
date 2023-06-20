package internal

import (
	"fmt"
)

func Scrape() {
	fmt.Println("==========     Scraping...     ==========")

	products, err := GetAllProducts()
	if err != nil {
		panic(err)
	}

	fmt.Println("==========     Scraping done!     ==========")
	fmt.Printf("Amount of products retrieved: %d\n", len(products))
	fmt.Println("==========     Saving to GCS ...     ==========")
	dbErr := SaveProductsToGCS(products)
	if dbErr != nil {
		panic(dbErr)
	}
	fmt.Println("==========     Saving to GCS done!     ==========")
}
