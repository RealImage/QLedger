package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func Dial() *sql.DB {
	conn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Unable to connect to Database:", err)
		panic(err)
	}
	log.Println("Successfully established connection to database.")
	return conn
}
