package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	shared "shared/pkg"
	"sort"
	"time"

	"cloud.google.com/go/storage"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	results, err := doQuery()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Got %d results\n", len(results))

	// Move existing pp.json to pp-old.json
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	from := client.Bucket(shared.GCSBucket).Object("prettige-prijzen/pp.json")
	to := client.Bucket(shared.GCSBucket).Object("prettige-prijzen/pp-old.json")
	if _, err := to.CopierFrom(from).Run(context.Background()); err != nil {
		panic(err)
	}
	fmt.Println("Moved pp.json to pp-old.json")
	// Delete pp.json
	if err := from.Delete(ctx); err != nil {
		panic(err)
	}
	fmt.Println("Deleted pp.json")

	// Save new results to pp.json
	doc := docObject{
		Date: time.Now(),
		Data: results,
	}
	serialized, err := json.Marshal(doc)
	if err != nil {
		panic(err)
	}
	err = shared.SaveJSONToGCS(shared.GCSBucket, "prettige-prijzen/pp.json", serialized)
	if err != nil {
		println("Issue writing to GCS: " + err.Error())
		panic(err)
	}

	// Save mini results to pp-mini.json
	ppMiniLength := 10
	if len(results) < ppMiniLength {
		ppMiniLength = len(results)
	}
	miniDoc := docObject{
		Date: time.Now(),
		Data: results[:ppMiniLength],
	}
	serialized, err = json.Marshal(miniDoc)
	if err != nil {
		panic(err)
	}
	err = shared.SaveJSONToGCS(shared.GCSBucket, "prettige-prijzen/pp-mini.json", serialized)
	if err != nil {
		println("Issue writing to GCS: " + err.Error())
		panic(err)
	}

	fmt.Println("Starting compare...")

	// Get pp-old.json to compare later
	obj := client.Bucket(shared.GCSBucket).Object("prettige-prijzen/pp-old.json")
	reader, err := obj.NewReader(ctx)
	if err != nil {
		panic(err)
	}
	defer reader.Close()
	var oldDoc docObject
	err = json.NewDecoder(reader).Decode(&oldDoc)
	if err != nil {
		panic(err)
	}
	// Parse to map for easy comparison
	oldMap := make(map[int]queryResult)
	for _, result := range oldDoc.Data {
		oldMap[result.ProductID] = result
	}

	// Compare new results to old results
	var newResults = []queryResult{}
	for _, result := range results {
		oldResult, ok := oldMap[result.ProductID]
		if !ok {
			newResults = append(newResults, result)
		} else if result.Diff-oldResult.Diff > 5 {
			newResults = append(newResults, result)
		}
	}
	fmt.Printf("There are %d new great deals!\n", len(newResults))
	// Save the compared results to pp-changes.json
	changesDoc := docObject{
		Date: time.Now(),
		Data: newResults,
	}
	serialized, err = json.Marshal(changesDoc)
	if err != nil {
		panic(err)
	}
	err = shared.SaveJSONToGCS(shared.GCSBucket, "prettige-prijzen/pp-changes.json", serialized)
	if err != nil {
		println("Issue writing to GCS: " + err.Error())
		panic(err)
	}
	// Save mini compared results to pp-changes-mini.json
	changesMiniLength := 10
	if len(newResults) < changesMiniLength {
		changesMiniLength = len(newResults)
	}
	changesMiniDoc := docObject{
		Date: time.Now(),
		Data: newResults[:changesMiniLength],
	}
	serialized, err = json.Marshal(changesMiniDoc)
	if err != nil {
		panic(err)
	}
	err = shared.SaveJSONToGCS(shared.GCSBucket, "prettige-prijzen/pp-changes-mini.json", serialized)
	if err != nil {
		println("Issue writing to GCS: " + err.Error())
		panic(err)
	}

	fmt.Println("PP Done, leggo!")
}

type queryResult struct {
	ProductID             int     `json:"productId"`
	LongName              string  `json:"longName"`
	SquareImage           string  `json:"squareImage"`
	BasicPrice            float64 `json:"basicPrice"`
	Benefit               string  `json:"benefit"`
	QuantityPrice         float64 `json:"quantityPrice"`
	QuantityPriceQuantity string  `json:"quantityPriceQuantity"`
	BestPrice             float64 `json:"bestPrice"`
	ThirtyDayAvg          float64 `json:"thirtyDayAvg"`
	Diff                  int     `json:"diff"`
}

type docObject struct {
	Date time.Time     `json:"date"`
	Data []queryResult `json:"data"`
}

func doQuery() (results []queryResult, err error) {
	connectionString := fmt.Sprintf(
		"postgres://%s?sslmode=disable",
		os.Getenv("DB_CONNECTION_STRING"),
	)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return
	}
	defer db.Close()
	db.SetMaxOpenConns(5)

	println("Executing query...")

	rows, err := db.Query(`
	SELECT computed_diff.*
	FROM (
		SELECT
			computed.*,
			((1 - (computed.best_price / computed.thirty_day_avg)) * 100)::int as diff
		FROM (
			SELECT DISTINCT
				pr.id,
				pr.long_name,
				pr.square_image,
				price.basic_price,
				COALESCE(promo.benefit, ''),
				price.quantity_price,
				price.quantity_price_quantity,
				(100 - (CASE
						WHEN promo.benefit IS NULL OR promo.benefit = '' -- 0%
						THEN 0::text
						WHEN SPLIT_PART(promo.benefit, ',', 2) != ''
						THEN SPLIT_PART(SPLIT_PART(promo.benefit, ',', 2), '_', 1)
						ELSE SPLIT_PART(promo.benefit, '_', 1)
					END)::numeric)
					* 0.01
					* (CASE WHEN price.quantity_price > 0
							THEN price.quantity_price
							ELSE price.basic_price
						END)
					as best_price,
				AVG(prices_for_avg.basic_price) as thirty_day_avg
			FROM
				products.product as pr
				INNER JOIN
				products.price as prices_for_avg
				ON pr.id = prices_for_avg.product_id
				INNER JOIN
				products.price as price
				ON pr.id = price.product_id
				LEFT JOIN
				products.promotion as promo
				ON promo.promotion_id = ANY(STRING_TO_ARRAY(price.promo_codes, ','))
			WHERE
				prices_for_avg.time > now() - interval '30 day'
				AND
				prices_for_avg.basic_price > 0
				AND
				pr.is_available = 'true'
				AND
				price.time > now() - interval '1 day'
				AND
				price.basic_price != 0
			GROUP BY
				pr.id,
				pr.square_image,
				pr.long_name,
				price.basic_price,
				promo.benefit,
				price.quantity_price,
				price.quantity_price_quantity
			HAVING
				AVG(prices_for_avg.basic_price) > 0
				AND
				(
					price.quantity_price != 0 -- quantity price
					OR
					promo.benefit IS NOT NULL -- active promotion
					OR
					price.basic_price < AVG(prices_for_avg.basic_price) * 0.97 -- normal price drop
				)
		) as computed
	) as computed_diff
	WHERE diff > 5
	ORDER BY diff DESC;
	`)
	if err != nil {
		return
	}
	defer rows.Close()

	println("Query executed.")
	resMap := make(map[int]queryResult)
	for rows.Next() {
		var tempResult queryResult
		err = rows.Scan(
			&tempResult.ProductID,
			&tempResult.LongName,
			&tempResult.SquareImage,
			&tempResult.BasicPrice,
			&tempResult.Benefit,
			&tempResult.QuantityPrice,
			&tempResult.QuantityPriceQuantity,
			&tempResult.BestPrice,
			&tempResult.ThirtyDayAvg,
			&tempResult.Diff,
		)
		if err != nil {
			println("lil issue over here: " + err.Error())
			continue
		}
		// Products can have multiple promotions at the same time
		if existVal, ok := resMap[tempResult.ProductID]; !ok {
			resMap[tempResult.ProductID] = tempResult
		} else {
			// Save the result with the highest diff
			if tempResult.Diff > existVal.Diff {
				resMap[tempResult.ProductID] = tempResult
			}
		}
	}
	err = rows.Err()
	if err = rows.Err(); err != nil {
		return []queryResult{}, err
	}
	for _, result := range resMap {
		results = append(results, result)
	}
	// Resort because we used a map
	sort.Slice(results, func(i, j int) bool {
		return results[i].Diff > results[j].Diff
	})
	return results, nil
}
