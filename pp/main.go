package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	shared "shared/pkg"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	connectionString := fmt.Sprintf(
		"postgres://%s?sslmode=disable",
		os.Getenv("DB_CONNECTION_STRING"),
	)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
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
						WHEN promo.benefit IS NULL -- 0%
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
				ON promo.promotion_id = price.promo_codes
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
		panic(err)
	}
	defer rows.Close()

	println("Query executed.")

	var results []queryResult = []queryResult{}
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
		results = append(results, tempResult)
	}

	fmt.Printf("Got %d results\n", len(results))

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

	miniDoc := docObject{
		Date: time.Now(),
		Data: results[:10],
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

	fmt.Println("PP Done!")
}

type queryResult struct {
	ProductID             int     `json:"product_id"`
	LongName              string  `json:"long_name"`
	SquareImage           string  `json:"square_image"`
	BasicPrice            float64 `json:"basic_price"`
	Benefit               string  `json:"benefit"`
	QuantityPrice         float64 `json:"quantity_price"`
	QuantityPriceQuantity string  `json:"quantity_price_quantity"`
	BestPrice             float64 `json:"best_price"`
	ThirtyDayAvg          float64 `json:"thirty_day_avg"`
	Diff                  int     `json:"diff"`
}

type docObject struct {
	Date time.Time     `json:"date"`
	Data []queryResult `json:"data"`
}
