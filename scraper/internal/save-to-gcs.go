package internal

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
	key := "colruyt-products/" + time.Now().Format("2006-01-02-15-04-05") + ".json"

	errWrite := shared.SaveJSONToGCS(shared.GCSBucket, key, serialized)
	if errWrite != nil {
		return errWrite
	}

	return nil
}
