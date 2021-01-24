package model

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
)

var (
	dbUser = os.Getenv("DB_USER")
	dbPw = os.Getenv("DB_PW")
	dbName = os.Getenv("DB_NAME")
	dbHost = os.Getenv("DB_HOST")
)

func ConnectDB() (*sql.DB, error) {
	//Note local DB 연결시 사용하는 코드
	db, err := sql.Open("postgres", "user=postgres password=1673" +
		" dbname=postgres host=localhost sslmode=disable port=5432")

	//Note 서버 DB 연결시 사용하는 코드
	//db, err := sql.Open("postgres",
	//	"user=redteam password=dkagh1234! dbname=redteam host=52.231.73.1 sslmode=disable port=5432")

	if db != nil {
		db.SetMaxOpenConns(100)
		db.SetMaxIdleConns(10)
	}
	if err != nil {
		return nil, err
	}
	return db, err
}
