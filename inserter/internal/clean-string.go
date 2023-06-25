package internal

import "strings"

func CleanString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
