package shared

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
)

func SaveJSONToGCS(
	bucket string,
	key string,
	serialized []byte,
) (
	err error,
) {

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	obj := client.Bucket(bucket).Object(key)
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
