package internal

import (
	"fmt"
	shared "shared/pkg"
	"strconv"
	"strings"
)

func GenerateInsertPromotionsQuery(
	promotions []shared.Promotion,
) string {

	query := `INSERT INTO products.promotion (
		` + strings.Join(shared.PromotionColumns, ", ") + `
	) VALUES `

	values := []string{}
	for _, promotion := range promotions {
		var benefit string
		var linkedProducts string
		for _, b := range promotion.Benefit {
			benefit += fmt.Sprintf("%v_%v_%v,", b.BenefitPercentage, b.MinLimit, b.LimitUnit)
		}
		benefit = strings.TrimSuffix(benefit, ",")
		for _, p := range promotion.LinkedProducts {
			linkedProducts += fmt.Sprintf("%d", p) + ","
		}
		linkedProducts = strings.TrimSuffix(linkedProducts, ",")

		v := `(
			'` + promotion.PromotionID + `',
			'` + promotion.ActiveStartDate + `',
			'` + promotion.ActiveEndDate + `',
			'` + benefit + `',
			'` + linkedProducts + `',
			'` + promotion.CommercialPromotionID + `',
			'` + promotion.FolderID + `',
			'` + fmt.Sprintf("%d", promotion.MaxTimes) + `',
			'` + strconv.FormatBool(promotion.Personalised) + `',
			'` + promotion.PromotionKind + `',
			'` + promotion.PromotionType + `',
			'` + promotion.PublicationStartDate + `',
			'` + promotion.PublicationEndDate + `'
		)`
		values = append(values, v)
	}

	query += strings.Join(values, ",")
	query += ` ON CONFLICT (promotion_id) DO UPDATE SET (
		` + strings.Join(shared.PromotionColumns, ",") + `
	) = (
		EXCLUDED.` + strings.Join(shared.PromotionColumns, ",EXCLUDED.") + `
	)`

	return query
}
