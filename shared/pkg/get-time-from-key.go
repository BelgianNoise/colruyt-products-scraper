package shared

import (
	"strings"
	"time"
)

func GetTimeFromKey(key string) (time.Time, error) {
	dateString := strings.Split(strings.Split(key, "/")[1], ".")[0]
	time, err := time.Parse("2006-01-02-15-04-05", dateString)
	if err != nil {
		return time, err
	}
	return time, nil
}
