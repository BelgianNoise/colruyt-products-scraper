package shared

import (
	"database/sql"
	"fmt"
	"os"
)

// Don't forget to call db.Close() after using the db instance
func CreateDBInstance() (*sql.DB, error) {
	connectionString := fmt.Sprintf(
		"postgres://%s?sslmode=disable",
		os.Getenv("DB_CONNECTION_STRING"),
	)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(5)

	return db, nil
}
