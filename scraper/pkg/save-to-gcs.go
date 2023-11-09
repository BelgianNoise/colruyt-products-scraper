package scraper

import (
	"encoding/json"
	"fmt"
	shared "shared/pkg"
	"time"
)

func SaveProductsToGCS(products []shared.Product) error {

	serialized, err := json.Marshal(products)
	if err != nil {
		return fmt.Errorf("failed to serialize products: %v", err)
	}
	key := "colruyt-products/" + time.Now().In(time.UTC).Format("2006-01-02-15-04-05") + ".json"

	errWrite := shared.SaveJSONToGCS(shared.GCSBucket, key, serialized)
	if errWrite != nil {
		return errWrite
	}

	return nil
}

func SavePromotionsToGCS(promotions []shared.Promotion) error {
	for _, promo := range promotions {
		serialized, err := json.Marshal(promo)
		if err != nil {
			return fmt.Errorf("failed to serialize promotion: %v", err)
		}
		key := "promotions/" + promo.PromotionID + ".json"

		errWrite := shared.SaveJSONToGCS(shared.GCSBucket, key, serialized)
		if errWrite != nil {
			return errWrite
		}
	}
	return nil
}
