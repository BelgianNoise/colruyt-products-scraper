package internal

import (
	"context"
	"fmt"
	shared "shared/pkg"
	"sort"
	"strings"
	"time"
)

func CompareTodayToPrevious(
	ctx context.Context,
	store *shared.FileStore,
) (
	jsonFileLocation string,
	e error,
) {
	objects, err := shared.GetProductsKeys(ctx, store)
	if err != nil {
		return "", err
	}
	if len(objects) == 0 {
		return "", fmt.Errorf("no product files found in storage")
	}
	latestKey := objects[0]
	today := time.Now().In(time.UTC).Format("2006-01-02")
	if !strings.HasPrefix(latestKey, shared.ProductsPrefix+today) {
		return "", fmt.Errorf("no file found for today, last file: %q, today: %q", latestKey, today)
	}

	lastFileTime, err := shared.GetTimeFromKey(latestKey)
	if err != nil {
		return "", err
	}

	// objects is ordered from most recent to oldest
	var previousFile string
	for _, file := range objects {
		fileTime, err := shared.GetTimeFromKey(file)
		if err != nil {
			continue
		}
		if fileTime.Before(lastFileTime.Add(-23 * time.Hour)) {
			previousFile = file
			break
		}
	}

	if previousFile == "" || latestKey == "" {
		return "", fmt.Errorf("no previous file %q or last file %q", previousFile, latestKey)
	}

	jsonFileLocation, changes, errCompare := shared.Compare(ctx, store, latestKey, previousFile, false)
	if errCompare != nil {
		return "", fmt.Errorf("error comparing %q to %q: %v", latestKey, previousFile, errCompare)
	}

	println("Sorting changes...")

	var increases = []shared.PriceDifference{}
	var decreases = []shared.PriceDifference{}
	for _, d := range changes {
		if d.PriceChangePercentage >= 0.01 {
			increases = append(increases, d)
		} else if d.PriceChangePercentage <= -0.01 {
			decreases = append(decreases, d)
		}
	}
	sort.Slice(increases, func(i, j int) bool {
		return increases[i].PriceChangePercentage > increases[j].PriceChangePercentage
	})
	sort.Slice(decreases, func(i, j int) bool {
		return decreases[i].PriceChangePercentage < decreases[j].PriceChangePercentage
	})

	println("Saving decreases...")
	err = store.PutJSON(ctx, "drastische-dalers/dd.json", map[string]interface{}{
		"date": time.Now(),
		"data": decreases,
	})
	if err != nil {
		return "", err
	}
	decMiniLength := 10
	if len(decreases) < decMiniLength {
		decMiniLength = len(decreases)
	}
	err = store.PutJSON(ctx, "drastische-dalers/dd-mini.json", map[string]interface{}{
		"date": time.Now(),
		"data": decreases[:decMiniLength],
	})
	if err != nil {
		return "", err
	}
	println("Done saving decreases!")

	println("Saving increases...")
	err = store.PutJSON(ctx, "sterkste-stijgers/ss.json", map[string]interface{}{
		"date": time.Now(),
		"data": increases,
	})
	if err != nil {
		return "", err
	}
	increasesMiniLength := 10
	if len(increases) < increasesMiniLength {
		increasesMiniLength = len(increases)
	}
	err = store.PutJSON(ctx, "sterkste-stijgers/ss-mini.json", map[string]interface{}{
		"date": time.Now(),
		"data": increases[:increasesMiniLength],
	})
	if err != nil {
		return "", err
	}
	println("Done saving increases!")

	return jsonFileLocation, nil
}
