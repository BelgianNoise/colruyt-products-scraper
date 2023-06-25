package internal

import (
	"database/sql"
	"fmt"
	"os"
	shared "shared/pkg"

	_ "github.com/lib/pq"
)

func InsertLatestData() error {
	key, err := shared.GetLatestProductsKey(shared.GCSBucket)
	if err != nil {
		return err
	}
	products, err := shared.GetProducts(shared.GCSBucket, key)
	if err != nil {
		return err
	}
	if len(products) == 0 {
		return nil
	}

	connectionString := fmt.Sprintf(
		"postgres://%s?sslmode=disable",
		os.Getenv("DB_CONNECTION_STRING"),
	)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}
	defer db.Close()
	db.SetMaxOpenConns(5)

	productChunks := [][]shared.Product{}
	chunkSize := 1000
	for i := 0; i < len(products); i += chunkSize {
		end := i + chunkSize
		if end > len(products) {
			end = len(products)
		}
		productChunks = append(productChunks, products[i:end])
	}

	fmt.Println("Starting insertion of", len(products), "products and prices")
	time, err := shared.GetTimeFromKey(key)
	if err != nil {
		return err
	}

	for _, chunk := range productChunks {

		prodsQuery := GenerateInsertProductsQuery(chunk)
		fmt.Println(" - Inserting or updating", len(chunk), "products")
		rows, err := db.Query(prodsQuery)
		if err != nil {
			return err
		}
		rows.Close()

		pricesQuery := GenerateInsertPricesQuery(chunk, time)
		fmt.Println(" - Inserting or updating", len(chunk), "prices")
		rows, err = db.Query(pricesQuery)
		if err != nil {
			return err
		}
		rows.Close()

	}

	fmt.Println("Insertion done!")
	return nil
}
