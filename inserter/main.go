package main

import (
	"inserter/internal"

	"github.com/joho/godotenv"
)

func main() {
	er := godotenv.Load(".env")
	if er != nil {
		panic(er)
	}
	err := internal.InsertLatestData()
	if err != nil {
		panic(err)
	}
}
