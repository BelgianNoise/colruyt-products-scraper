package shared

import (
	"encoding/json"
	"sort"
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
	sort.Sort(sort.Reverse(sort.StringSlice(allobjs)))

	return allobjs[0], nil
}
