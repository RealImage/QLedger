package database

import (
	"log"
	"os"
	"sync"

	"gopkg.in/jackc/pgx.v2"
)

var ConnPool *pgx.ConnPool
var once sync.Once

func init() {
	if ConnPool == nil {
		once.Do(func() {
			connConfig, err := pgx.ParseURI(os.Getenv("DATABASE_URL"))
			if err != nil {
				log.Fatal("Invalid Database URL:", err)
				panic(err)
				return
			}
			poolConfig := pgx.ConnPoolConfig{
				ConnConfig:     connConfig,
				MaxConnections: 10, //TODO: Read from config
			}
			ConnPool, err = pgx.NewConnPool(poolConfig)
			if err != nil {
				log.Fatal("Unable to connect to Database:", err)
				panic(err)
				return
			}
			log.Println("Successfully established connection to database:", poolConfig.Database)
		})
	}
}

func Cleanup() {
	if ConnPool == nil {
		ConnPool.Close()
	}
}
