package shared

import (
	"encoding/json"
	"fmt"
	"strings"
)

func Compare(
	laterFile string,
	earlierFile string,
	excludePromotions bool,
) (
	jsonFileLocation string,
	err error,
) {
	fmt.Printf("Comparing %q to %q\n", laterFile, earlierFile)

	key := fmt.Sprintf(
		"price-changes/%v_%v/promotions-%v.json",
		strings.Split(strings.Split(earlierFile, "/")[1], ".")[0],
		strings.Split(strings.Split(laterFile, "/")[1], ".")[0],
		!excludePromotions,
	)

	d, _ := GetObjectFromBucket(GCSBucket, key)
	if len(d) > 0 {
		fmt.Printf("Already compared %q to %q\n", laterFile, earlierFile)
		return key, nil
	}

	laterList, err := GetProducts(GCSBucket, laterFile)
	if err != nil {
		return "", err
	}
	earlierList, err := GetProducts(GCSBucket, earlierFile)
	if err != nil {
		return "", err
	}

	diff := []PriceDifference{}
	up := 0
	down := 0

	for _, laterProduct := range laterList {
		for _, earlierProduct := range earlierList {
			if laterProduct.ProductID == earlierProduct.ProductID {
				// Don't include promotions
				if excludePromotions && (laterProduct.Price.IsRedPrice || earlierProduct.Price.IsRedPrice) {
					continue
				}
				if laterProduct.Price.BasicPrice != 0 && earlierProduct.Price.BasicPrice != 0 {
					if laterProduct.Price.BasicPrice != earlierProduct.Price.BasicPrice {
						change := laterProduct.Price.BasicPrice - earlierProduct.Price.BasicPrice
						diff = append(diff, PriceDifference{
							LongName:              laterProduct.LongName,
							PriceChange:           change,
							PriceChangePercentage: (change) / earlierProduct.Price.BasicPrice,
							InvolvesPromotion:     len(laterProduct.Promotion) > 0 || len(earlierProduct.Promotion) > 0 || laterProduct.Price.IsRedPrice || earlierProduct.Price.IsRedPrice,
							OldPrice:              earlierProduct.Price,
							Price:                 laterProduct.Price,
							Product:               laterProduct,
						})
						if change > 0 {
							up++
						} else {
							down++
						}
					}
				}
				// product found, no need to keep looping
				break
			}
		}
	}

	fmt.Printf("%d products have changed prices (up: %d | down: %d)\n", len(diff), up, down)

	serialized, err := json.Marshal(diff)
	if err != nil {
		return "", err
	}

	saveErr := SaveJSONToGCS(GCSBucket, key, serialized)
	if saveErr != nil {
		return "", saveErr
	}

	return key, nil
}
