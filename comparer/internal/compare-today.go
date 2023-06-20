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
	objects, err := shared.ListBucketObjects(shared.GCSBucket)
	if err != nil {
		return "", err
	}
	sort.Sort(sort.Reverse(sort.StringSlice(objects)))

	lastFile := objects[0]

	today := time.Now().In(time.UTC).Format("2006-01-02")
	if !strings.HasPrefix(lastFile, shared.GCSBucket+"/"+today) {
		return "", fmt.Errorf("no file found for today, last file: %q, today: %q", lastFile, today)
	}

	var previousFile string
	for _, file := range objects {
		dateString := strings.Split(strings.Split(file, "/")[1], ".")[0]
		fileTime, err := time.Parse("2006-01-02-15-04-05", dateString)
		if err != nil {
			continue
		}
		if fileTime.Before(time.Now().Add(-20 * time.Hour)) {
			previousFile = file
			break
		}
	}

	if previousFile == "" || lastFile == "" {
		return "", fmt.Errorf("no previous file %q or last file %q", previousFile, lastFile)
	}

	jsonFileLocation, errCompare := shared.Compare(lastFile, previousFile)
	if errCompare != nil {
		return "", errCompare
	}

	return jsonFileLocation, nil
}
