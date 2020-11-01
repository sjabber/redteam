package model

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("postgres",
		"user=postgres password=1673 dbname=postgres host=localhost sslmode=disable")
	if db != nil {
		db.SetMaxOpenConns(100)
		db.SetMaxIdleConns(10)
	}
	if err != nil {
		return nil, err
	}
	return db, err
}