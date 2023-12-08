package internal

import (
	"encoding/json"
	"fmt"
	shared "shared/pkg"
	"sort"
	"strings"
	"time"
)

func CompareTodayToPrevious() (
	jsonFileLocation string,
	e error,
) {
	latestKey, err := shared.GetLatestProductsKey(shared.GCSBucket)
	if err != nil {
		return "", err
	}
	today := time.Now().In(time.UTC).Format("2006-01-02")
	if !strings.HasPrefix(latestKey, shared.GCSBucket+"/"+today) {
		return "", fmt.Errorf("no file found for today, last file: %q, today: %q", latestKey, today)
	}

	lastFileTime, err := shared.GetTimeFromKey(latestKey)
	if err != nil {
		return "", err
	}

	objects, err := shared.ListBucketObjects(shared.GCSBucket, "colruyt-products/")
	if err != nil {
		return "", err
	}
	// Look from most recent to oldest
	sort.Sort(sort.Reverse(sort.StringSlice(objects)))

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

	jsonFileLocation, changes, errCompare := shared.Compare(latestKey, previousFile, false)
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
	decreasesSerialized, err := json.Marshal(map[string]interface{}{
		"date": time.Now(),
		"data": decreases,
	})
	if err != nil {
		return "", err
	}
	err = shared.SaveJSONToGCS(shared.GCSBucket, "drastische-dalers/dd.json", decreasesSerialized)
	if err != nil {
		return "", err
	}
	decMiniLength := 10
	if len(decreases) < decMiniLength {
		decMiniLength = len(decreases)
	}
	decreasesMiniSerialized, err := json.Marshal(map[string]interface{}{
		"date": time.Now(),
		"data": decreases[:decMiniLength],
	})
	if err != nil {
		return "", err
	}
	err = shared.SaveJSONToGCS(shared.GCSBucket, "drastische-dalers/dd-mini.json", decreasesMiniSerialized)
	if err != nil {
		return "", err
	}
	println("Done saving decreases!")

	println("Saving increases...")
	increasesSerialized, err := json.Marshal(map[string]interface{}{
		"date": time.Now(),
		"data": increases,
	})
	if err != nil {
		return "", err
	}
	err = shared.SaveJSONToGCS(shared.GCSBucket, "sterkste-stijgers/ss.json", increasesSerialized)
	if err != nil {
		return "", err
	}
	increasesMiniLength := 10
	if len(increases) < increasesMiniLength {
		increasesMiniLength = len(increases)
	}
	increasesMiniSerialized, err := json.Marshal(map[string]interface{}{
		"date": time.Now(),
		"data": increases[:increasesMiniLength],
	})
	if err != nil {
		return "", err
	}
	err = shared.SaveJSONToGCS(shared.GCSBucket, "sterkste-stijgers/ss-mini.json", increasesMiniSerialized)
	if err != nil {
		return "", err
	}
	println("Done saving increases!")

	return jsonFileLocation, nil
}
