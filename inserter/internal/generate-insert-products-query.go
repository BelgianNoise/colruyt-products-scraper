package internal

import (
	shared "shared/pkg"
	"strconv"
	"strings"
)

func GenerateInsertProductsQuery(
	products []shared.Product,
) string {
	query := `INSERT INTO products.product (
		` + strings.Join(shared.ProductColumns, ", ") + `
	) VALUES `

	values := []string{}
	for _, product := range products {
		v := `(
			'` + CleanString(product.ProductID) + `',
			'` + CleanString(product.Name) + `',
			'` + CleanString(product.LongName) + `',
			'` + CleanString(product.ShortName) + `',
			'` + CleanString(product.Content) + `',
			'` + CleanString(product.FullImage) + `',
			'` + CleanString(product.SquareImage) + `',
			'` + CleanString(product.ThumbNail) + `',
			'` + CleanString(product.CommercialArticleNumber) + `',
			'` + CleanString(product.TechnicalArticleNumber) + `',
			'` + CleanString(product.AlcoholVolume) + `',
			'` + CleanString(product.CountryOfOrigin) + `',
			'` + CleanString(product.FicCode) + `',
			'` + strconv.FormatBool(product.IsBiffe) + `',
			'` + strconv.FormatBool(product.IsBio) + `',
			'` + strconv.FormatBool(product.IsExclusivelySoldInLuxembourg) + `',
			'` + strconv.FormatBool(product.IsNew) + `',
			'` + strconv.FormatBool(product.IsPrivateLabel) + `',
			'` + strconv.FormatBool(product.IsWeightArticle) + `',
			'` + CleanString(product.OrderUnit) + `',
			'` + CleanString(product.RecentQuanityOfStockUnits) + `',
			'` + CleanString(product.WeightconversionFactor) + `',
			'` + CleanString(product.Brand) + `',
			'` + CleanString(product.BusinessDomain) + `',
			'` + strconv.FormatBool(product.IsAvailable) + `',
			'` + CleanString(product.SeoBrand) + `',
			'` + CleanString(product.TopCategoryId) + `',
			'` + CleanString(product.TopCategoryName) + `',
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

	return query
}
