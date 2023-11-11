package internal

import (
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

	jsonFileLocation, errCompare := shared.Compare(latestKey, previousFile, false)
	if errCompare != nil {
		return "", fmt.Errorf("error comparing %q to %q: %v", latestKey, previousFile, errCompare)
	}

	return jsonFileLocation, nil
}
