package main

import (
	"encoding/json"
	"fmt"
	"os"
	scraper "scraper/pkg"
	shared "shared/pkg"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	scraper.ScraperAPIKey = os.Getenv("SCRAPER_API_KEY")
	// r, err := scraper.DoAPICall(1, 250)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(r.ProductsReturned)
	objs, err := shared.ListBucketObjects(shared.GCSBucket)
	if err != nil {
		panic(err)
	}
	for _, obj := range objs {
		if !strings.Contains(obj, "colruyt-products") {
			continue
		}
		f, err := shared.GetObjectFromBucket(shared.GCSBucket, obj)
		if err != nil {
			panic(err)
		}
		var prods []shared.Product
		err = json.Unmarshal(f, &prods)
		if err != nil {
			panic(err)
		}

		fmt.Println(len(prods))
	}
}
