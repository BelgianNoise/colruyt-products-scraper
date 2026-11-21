package internal

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	shared "shared/pkg"
	"time"

	_ "github.com/lib/pq"
)

func InsertLatestData(ctx context.Context, store *shared.FileStore) error {
	db, dbError := shared.CreateDBInstance()
	if dbError != nil {
		return dbError
	}
	defer db.Close()

	err := InsertProductsAndPrices(ctx, store, db)
	if err != nil {
		return err
	}

	fmt.Println("Upserting promotions ...")
	err = UpsertPromotionData(ctx, store, db)
	if err != nil {
		return err
	}
	fmt.Println("Upserting promotions done!")
	return nil
}

func InsertProductsAndPrices(
	ctx context.Context,
	store *shared.FileStore,
	db *sql.DB,
) error {
	key, err := shared.GetLatestProductsKey(ctx, store)
	if err != nil {
		return err
	}
	fmt.Println("==========      Inserting data from", key)
	products, err := shared.GetProducts(ctx, store, key)
	if err != nil {
		return err
	}
	if len(products) == 0 {
		return nil
	}

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
	fmt.Println("Inserting products and prices done!")

	return nil
}

// Why we upsert ? : https://github.com/BelgianNoise/colruyt-products-scraper/issues/15
func UpsertPromotionData(
	ctx context.Context,
	store *shared.FileStore,
	db *sql.DB,
) error {
	files, err := store.List(ctx, "promotions/")
	if err != nil {
		return err
	}
	// Promotions are rewritten on every scrape, so recently modified
	// files are the ones still active.
	since := time.Now().Add(-48 * time.Hour)
	var promotionObjects []string
	for _, f := range files {
		if f.LastModified.After(since) {
			promotionObjects = append(promotionObjects, f.Key)
		}
	}

	fmt.Println("Found", len(promotionObjects), "promotions in storage to upsert")
	var promotions []shared.Promotion
	for _, object := range promotionObjects {
		data, err := store.Get(ctx, object)
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

	if len(promotions) == 0 {
		fmt.Println("No new promotions to upsert")
		return nil
	} else {
		fmt.Println("Upserting", len(promotions), "promotions")
		queryStr := GenerateInsertPromotionsQuery(promotions)
		_, err = db.Query(queryStr)
		if err != nil {
			return err
		}
	}

	return nil
}
