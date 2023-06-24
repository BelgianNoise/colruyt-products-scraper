package internal

import (
	"database/sql"
	"fmt"
	"os"
	shared "shared/pkg"
	"strconv"
	"strings"

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

	fmt.Println("Starting insertion of", len(products), "products")

	for _, chunk := range productChunks {

		query := `INSERT INTO products.product (
			` + strings.Join(shared.ProductColumns, ", ") + `
		) VALUES `

		values := []string{}
		for _, product := range chunk {
			v := `(
				'` + cleanString(product.ProductID) + `',
				'` + cleanString(product.Name) + `',
				'` + cleanString(product.LongName) + `',
				'` + cleanString(product.ShortName) + `',
				'` + cleanString(product.Content) + `',
				'` + cleanString(product.FullImage) + `',
				'` + cleanString(product.SquareImage) + `',
				'` + cleanString(product.ThumbNail) + `',
				'` + cleanString(product.CommercialArticleNumber) + `',
				'` + cleanString(product.TechnicalArticleNumber) + `',
				'` + cleanString(product.AlcoholVolume) + `',
				'` + cleanString(product.CountryOfOrigin) + `',
				'` + cleanString(product.FicCode) + `',
				'` + strconv.FormatBool(product.IsBiffe) + `',
				'` + strconv.FormatBool(product.IsBio) + `',
				'` + strconv.FormatBool(product.IsExclusivelySoldInLuxembourg) + `',
				'` + strconv.FormatBool(product.IsNew) + `',
				'` + strconv.FormatBool(product.IsPrivateLabel) + `',
				'` + strconv.FormatBool(product.IsWeightArticle) + `',
				'` + cleanString(product.OrderUnit) + `',
				'` + cleanString(product.RecentQuanityOfStockUnits) + `',
				'` + cleanString(product.WeightconversionFactor) + `',
				'` + cleanString(product.Brand) + `',
				'` + cleanString(product.BusinessDomain) + `',
				'` + strconv.FormatBool(product.IsAvailable) + `',
				'` + cleanString(product.SeoBrand) + `',
				'` + cleanString(product.TopCategoryId) + `',
				'` + cleanString(product.TopCategoryName) + `',
				'` + strconv.Itoa(product.WalkRouteSequenceNumber) + `'
			)`
			values = append(values, v)
		}

		query += strings.Join(values, ",")
		query += ` ON CONFLICT (id) DO UPDATE SET (
			` + strings.Join(shared.ProductColumns, ",") + `
		) = (
			EXCLUDED.` + strings.Join(shared.ProductColumns, ",EXCLUDED.") + `
		)`

		fmt.Println(" - Inserting or updating", len(chunk), "products")
		// fmt.Println(query)
		rows, err := db.Query(query)
		if err != nil {
			return err
		}
		rows.Close()
	}

	fmt.Println("Insertion done!")
	return nil
}

func cleanString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
