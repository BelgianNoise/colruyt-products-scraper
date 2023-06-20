package main

import (
	internal "comparer/internal"
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	fmt.Printf("==========      Comparing      ============\n")
	file, errCompare := internal.CompareTodayToPrevious()
	if errCompare != nil {
		panic(errCompare)
	}
	fmt.Printf("======      Compareing done! File location: %q     ======\n", file)
}
