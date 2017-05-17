package database

import (
	"log"
	"os"
	"sync"
	"database/sql"

	_ "github.com/lib/pq"
)

var Conn *sql.DB
var once sync.Once

func init() {
	if Conn == nil {
		once.Do(func() {
			var err error
			Conn, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
			if err != nil {
				log.Fatal("Unable to connect to Database:", err)
				panic(err)
				return
			}
			log.Println("Successfully established connection to database.")
		})
	}
}

func Cleanup() {
	if Conn == nil {
		err := Conn.Close()
		if err != nil {
			log.Fatal("Error closing db connection:", err)
		}
	}
}
