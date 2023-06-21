package shared

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func GetObjectFromBucket(
	bucket string,
	key string,
) (
	data []byte,
	e error,
) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithoutAuthentication())
	if err != nil {
		return []byte{}, err
	}
	defer client.Close()

	rc, err := client.Bucket(bucket).Object(key).NewReader(ctx)
	if err != nil {
		return []byte{}, err
	}
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil {
		return []byte{}, err
	}

	return content, nil
}
