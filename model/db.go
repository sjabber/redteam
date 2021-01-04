package model

import (
	"database/sql"
	_ "github.com/lib/pq"
)

//var (
//	dbUser = os.Getenv("DB_USER")
//	dbPw = os.Getenv("DB_PW")
//	dbName = os.Getenv("DB_NAME")
//	dbHost = os.Getenv("DB_HOST")
//)

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("postgres",
		"user=redteam password=dkagh1234! dbname=redteam host=52.231.73.1 sslmode=disable port=5432")
	if db != nil {
		db.SetMaxOpenConns(100)
		db.SetMaxIdleConns(10)
	}
	if err != nil {
		return nil, err
	}
	return db, err
}
