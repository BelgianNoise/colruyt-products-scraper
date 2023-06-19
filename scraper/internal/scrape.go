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
	fmt.Println("==========     Saving to R2 DB...     ==========")
	dbErr := SaveToGCS([]Product{})
	if dbErr != nil {
		panic(dbErr)
	}
	fmt.Println("==========     Saving to R2 DB done!     ==========")
}
