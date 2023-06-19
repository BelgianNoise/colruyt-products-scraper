package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
)

func SaveToGCS(products []Product) error {

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	serialized, err := json.Marshal(products)
	if err != nil {
		return fmt.Errorf("failed to serialize products: %v", err)
	}
	key := "colruyt-products/" + time.Now().Format("2006-01-02-15-04-05") + ".json"

	obj := client.Bucket(GCSBucket).Object(key)
	writer := obj.NewWriter(ctx)

	if _, err := writer.Write(serialized); err != nil {
		return fmt.Errorf("failed to write to bucket: %v", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	// Set Content-Type to application/json
	attrs := storage.ObjectAttrsToUpdate{ContentType: "application/json"}
	if _, err := obj.Update(ctx, attrs); err != nil {
		return fmt.Errorf("failed trying to set metadata: %s", err)
	}

	return nil
}
