package internal

import (
	"fmt"
	shared "shared/pkg"
	"strconv"
	"strings"
	"time"
)

func GenerateInsertPricesQuery(
	products []shared.Product,
	timeString time.Time,
) string {

	query := `INSERT INTO products.price (
		` + strings.Join(shared.PriceColumns, ", ") + `
	) VALUES `

	values := []string{}
	for _, product := range products {
		promos := []string{}
		for _, promo := range product.Promotion {
			promos = append(promos, fmt.Sprintf("%s_%s_%s", promo.ID, promo.PublicationStartDate, promo.PublicationEndDate))
		}
		promosString := "NULL"
		if len(promos) > 0 {
			promosString = strings.Join(promos, ",")
		}
		v := `(
			'` + CleanString(product.ProductID) + `',
			'` + fmt.Sprintf("%f", product.Price.BasicPrice) + `',
			'` + strconv.FormatBool(product.Price.IsRedPrice) + `',
			'` + strconv.FormatBool(product.InPromo) + `',
			'` + strconv.FormatBool(product.InPreConditionPromo) + `',
			'` + strconv.FormatBool(product.IsPriceAvailable) + `',
			'` + CleanString(product.Price.MeasurementUnit) + `',
			'` + fmt.Sprintf("%f", product.Price.MeasurementUnitPrice) + `',
			'` + CleanString(product.Price.RecommendedQuantity) + `',
			'` + timeString.UTC().Format(time.RFC3339) + `',
			'` + promosString + `'
		)`
		values = append(values, v)
	}

	query += strings.Join(values, ",")
	query += ` ON CONFLICT DO NOTHING`

	return query
}
