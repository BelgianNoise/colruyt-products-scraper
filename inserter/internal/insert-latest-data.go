package internal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	shared "shared/pkg"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func InsertLatestData() error {
	key, err := shared.GetLatestProductsKey(shared.GCSBucket)
	if err != nil {
		return err
	}
	fmt.Println("==========      Inserting data from", key)
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

	fmt.Println("Inserting promotion ...")
	err = InsertPromotionData(db)
	if err != nil {
		return err
	}
	fmt.Println("Insertion done!")
	return nil
}

func InsertPromotionData(
	db *sql.DB,
) error {
	promotionObjects, err := shared.ListBucketObjectsInTimeRange(
		shared.GCSBucket,
		"promotions/",
		time.Now().Add(20*time.Hour*-1),
		time.Now(),
	)
	if err != nil {
		return err
	}

	var promotions []shared.Promotion
	for _, object := range promotionObjects {
		promotionID := strings.Split(strings.Split(object, "/")[1], ".")[0]
		exists := false
		r := db.QueryRow(
			"SELECT promotion_id FROM products.promotion WHERE promotion_id = $1",
			promotionID,
		)
		err := r.Scan(&promotionID)
		if err == sql.ErrNoRows {
			exists = false
		} else if err != nil {
			return err
		} else {
			exists = true
		}

		if exists {
			fmt.Println("Promotion already in DB: ", object)
			continue
		} else {
			fmt.Println("New Promotion Found: ", object)
			data, err := shared.GetObjectFromBucket(shared.GCSBucket, object)
			if err != nil {
				return err
			}
			var promotion shared.Promotion
			err = json.Unmarshal(data, &promotion)
			if err != nil {
				return err
			}
			promotions = append(promotions, promotion)
		}
	}

	if len(promotions) == 0 {
		fmt.Println("No new promotions to insert")
		return nil
	} else {
		queryStr := GenerateInsertPromotionsQuery(promotions)
		_, err = db.Query(queryStr)
		if err != nil {
			return err
		}
	}

	return nil
}
