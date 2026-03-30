package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func OpenDB() *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("error: ", err)
	}
	if err := db.Ping(); err != nil {
		log.Println("error: ", err)
	}
	return db
}
