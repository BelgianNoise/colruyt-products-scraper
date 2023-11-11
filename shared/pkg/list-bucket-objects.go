package shared

import (
	"context"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func ListBucketObjects(
	bucket string,
	prefix string,
) (
	objects []string,
	err error,
) {

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithoutAuthentication())
	if err != nil {
		return []string{}, err
	}

	defer client.Close()

	buck := client.Bucket(bucket)
	// list all objects in the colruyt-product folder
	objs := buck.Objects(ctx, &storage.Query{
		Prefix: prefix,
	})

	for {
		attrs, err := objs.Next()
		if err == iterator.Done {
			break
		}
		if err == nil {
			objects = append(objects, attrs.Name)
		}
	}

	return objects, nil
}
