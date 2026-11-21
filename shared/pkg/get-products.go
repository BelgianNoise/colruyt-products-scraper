package shared

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
)

// ProductsPrefix is the folder that holds the daily product dumps.
const ProductsPrefix = "colruyt-products/"

func GetProducts(
	ctx context.Context,
	store *FileStore,
	key string,
) (
	products []Product,
	err error,
) {
	data, err := store.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func GetLatestProducts(
	ctx context.Context,
	store *FileStore,
) (
	products []Product,
	err error,
) {
	latestKey, err := GetLatestProductsKey(ctx, store)
	if err != nil {
		return nil, err
	}
	return GetProducts(ctx, store, latestKey)
}

// GetProductsKeys returns all product dump keys, newest first.
func GetProductsKeys(
	ctx context.Context,
	store *FileStore,
) (
	keys []string,
	err error,
) {
	files, err := store.List(ctx, ProductsPrefix)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		keys = append(keys, f.Key)
	}
	// Keys contain a sortable timestamp, so reverse order is newest first.
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	return keys, nil
}

func GetLatestProductsKey(
	ctx context.Context,
	store *FileStore,
) (
	key string,
	err error,
) {
	keys, err := GetProductsKeys(ctx, store)
	if err != nil {
		return "", err
	}
	if len(keys) == 0 {
		return "", errors.New("no product files found in storage")
	}
	return keys[0], nil
}
