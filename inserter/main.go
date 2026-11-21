package main

import (
	"context"
	"inserter/internal"
	shared "shared/pkg"

	"github.com/joho/godotenv"
)

func main() {
	er := godotenv.Load(".env")
	if er != nil {
		panic(er)
	}
	store, err := shared.NewFileStore()
	if err != nil {
		panic(err)
	}
	err = internal.InsertLatestData(context.Background(), store)
	if err != nil {
		panic(err)
	}
}
