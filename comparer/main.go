package main

import (
	internal "comparer/internal"
	"context"
	"fmt"
	shared "shared/pkg"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	store, err := shared.NewFileStore()
	if err != nil {
		panic(err)
	}
	fmt.Printf("==========      Comparing      ============\n")
	file, errCompare := internal.CompareTodayToPrevious(context.Background(), store)
	if errCompare != nil {
		panic(errCompare)
	}
	fmt.Printf("======      Compareing done! File location: %q     ======\n", file)
}
