package internal

import (
	"fmt"
	shared "shared/pkg"
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

	lastFileDateString := strings.Split(strings.Split(latestKey, "/")[1], ".")[0]
	lastFileTime, err := time.Parse("2006-01-02-15-04-05", lastFileDateString)
	if err != nil {
		return "", err
	}

	objects, err := shared.ListBucketObjects(shared.GCSBucket)
	if err != nil {
		return "", err
	}
	var previousFile string
	for _, file := range objects {
		dateString := strings.Split(strings.Split(file, "/")[1], ".")[0]
		fileTime, err := time.Parse("2006-01-02-15-04-05", dateString)
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
