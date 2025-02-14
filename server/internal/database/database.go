package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(uri string) {
	var err error
	DB, err = sql.Open("postgres", uri)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if DB == nil {
		log.Fatal("database connection is nil")
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	log.Println("Connected to the PostgreSQL database successfully")

}

func CloseDB() {
	err := DB.Close()
	if err != nil {
		log.Printf("Error closing the database: %v", err)
	}
}
