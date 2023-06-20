package main

import (
	internal "comparer/internal"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	internal.Compare("", "")
}
