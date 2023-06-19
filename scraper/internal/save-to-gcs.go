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
	fmt.Println("11111")
	client, err := storage.NewClient(ctx)
	fmt.Println("22222")
	if err != nil {
		fmt.Println("33333")
		return err
	}
	fmt.Println("44444")

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

	return nil
}
