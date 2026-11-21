package shared

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func Compare(
	ctx context.Context,
	store *FileStore,
	laterFile string,
	earlierFile string,
	excludePromotions bool,
) (
	jsonFileLocation string,
	diff []PriceDifference,
	err error,
) {
	fmt.Printf("Comparing %q to %q\n", laterFile, earlierFile)

	key := fmt.Sprintf(
		"price-changes/%v_%v/promotions-%v.json",
		strings.Split(strings.Split(earlierFile, "/")[1], ".")[0],
		strings.Split(strings.Split(laterFile, "/")[1], ".")[0],
		!excludePromotions,
	)

	d, err := store.Get(ctx, key)
	if err == nil {
		fmt.Printf("Already compared %q to %q\n", laterFile, earlierFile)
		err = json.Unmarshal(d, &diff)
		if err != nil {
			return "", diff, err
		}
		return key, diff, nil
	}
	if !errors.Is(err, ErrNotFound) {
		return "", diff, err
	}

	laterList, err := GetProducts(ctx, store, laterFile)
	if err != nil {
		return "", diff, err
	}
	earlierList, err := GetProducts(ctx, store, earlierFile)
	if err != nil {
		return "", diff, err
	}

	diff, up, down := CompareLists(laterList, earlierList, excludePromotions)

	fmt.Printf("%d products have changed prices (up: %d | down: %d)\n", len(diff), up, down)

	if err := store.PutJSON(ctx, key, diff); err != nil {
		return "", diff, err
	}

	return key, diff, nil
}

func CompareLists(
	laterList []Product,
	earlierList []Product,
	excludePromotions bool,
) (
	diff []PriceDifference,
	up int,
	down int,
) {
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

	return diff, up, down
}
