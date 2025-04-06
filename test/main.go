package main

import (
	"fmt"
	scraper "scraper/pkg"
	shared "shared/pkg"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	scraper.InitVariables()
	scraper.LoadCookies()
	shared.InitProxyVars()

	products, err := scraper.GetAllProductsWithParams(30.0, 2, 250, false)
	if err != nil {
		panic(err)
	}

	for _, product := range products {
		if len(product.Promotion) > 0 {
			println(product.Promotion[0].TechPromoID)
		}
	}
	promotions, err := scraper.GetAllPromotions(products)
	if err != nil {
		panic(err)
	}
	println(len(promotions))
	fmt.Printf("%v", promotions[0])
}
