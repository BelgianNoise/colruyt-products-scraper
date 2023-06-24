package shared

import (
	"encoding/json"
	"sort"
	"strings"
)

func GetProducts(
	bucket string,
	file string,
) (
	products []Product,
	err error,
) {
	data, err := GetObjectFromBucket(bucket, file)
	if err != nil {
		return nil, err
	}
	errJSONLater := json.Unmarshal(data, &products)
	if err != nil {
		return nil, errJSONLater
	}
	return products, nil
}

func GetLatestProducts(
	bucket string,
) (
	products []Product,
	err error,
) {
	latestKey, err := GetLatestProductsKey(bucket)
	if err != nil {
		return nil, err
	}
	return GetProducts(bucket, latestKey)
}

func GetLatestProductsKey(
	bucket string,
) (
	key string,
	err error,
) {
	allobjs, err := ListBucketObjects(bucket)
	if err != nil {
		return "", err
	}
	objs := []string{}
	for _, obj := range allobjs {
		if strings.HasPrefix(obj, "colruyt-products/") {
			objs = append(objs, obj)
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(objs)))

	return objs[0], nil
}
