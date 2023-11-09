package main

import (
	"fmt"
	scraper "scraper/pkg"
	shared "shared/pkg"
)

func main() {
	prods, err := shared.GetLatestProducts(shared.GCSBucket)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Amount of products retrieved: %d\n", len(prods))
	promos := scraper.GetAllPromotions(prods)
	fmt.Println("Got promos: ", len(promos))
}
