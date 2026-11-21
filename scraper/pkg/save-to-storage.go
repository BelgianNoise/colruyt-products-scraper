package scraper

import (
	"context"
	shared "shared/pkg"
	"time"
)

func SaveProducts(ctx context.Context, store *shared.FileStore, products []shared.Product) error {
	key := shared.ProductsPrefix + time.Now().In(time.UTC).Format("2006-01-02-15-04-05") + ".json"
	return store.PutJSON(ctx, key, products)
}

func SavePromotions(ctx context.Context, store *shared.FileStore, promotions []shared.Promotion) error {
	for _, promo := range promotions {
		key := "promotions/" + promo.PromotionID + ".json"
		if err := store.PutJSON(ctx, key, promo); err != nil {
			return err
		}
	}
	return nil
}
