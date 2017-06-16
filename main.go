package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	ledgerContext "github.com/RealImage/QLedger/context"
	"github.com/RealImage/QLedger/controllers"
	"github.com/RealImage/QLedger/middlewares"
	"github.com/julienschmidt/httprouter"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Panic("Unable to connect to Database:", err)
	}
	log.Println("Successfully established connection to database.")

	log.Println("Starting db schema migration...")
	driver, _ := postgres.WithInstance(db, &postgres.Config{})
	m, _ := migrate.NewWithDatabaseInstance(
		"file://migrations/postgres",
		"postgres", driver)
	version, _, _ := m.Version()
	log.Println("Current schema version:", version)
	err = m.Up()
	if err != nil {
		log.Println("Error while migration:", err)
	}
	version, _, _ = m.Version()
	log.Println("Migrated schema version:", version)
	appContext := &ledgerContext.AppContext{DB: db}

	router := httprouter.New()

	// Create accounts and transactions
	// router.HandlerFunc("POST", "/v1/accounts", middlewares.ContextMiddleware(controllers.AddAccount, appContext))
	router.HandlerFunc("POST", "/v1/transactions", middlewares.ContextMiddleware(controllers.MakeTransaction, appContext))

	// Read or search accounts and transactions
	router.HandlerFunc("GET", "/v1/accounts", middlewares.ContextMiddleware(controllers.GetAccounts, appContext))
	router.HandlerFunc("POST", "/v1/accounts/_search", middlewares.ContextMiddleware(controllers.GetAccounts, appContext))
	router.HandlerFunc("GET", "/v1/transactions", middlewares.ContextMiddleware(controllers.GetTransactions, appContext))
	router.HandlerFunc("POST", "/v1/transactions/_search", middlewares.ContextMiddleware(controllers.GetTransactions, appContext))

	// Update data of accounts and transactions
	// router.HandlerFunc("PUT", "/v1/accounts", middlewares.ContextMiddleware(controllers.UpdateAccount, appContext))
	// router.HandlerFunc("PUT", "/v1/transactions", middlewares.ContextMiddleware(controllers.UpdateTransaction, appContext))

	port := "7000"
	if lp := os.Getenv("PORT"); lp != "" {
		port = lp
	}
	log.Println("Running server on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, router))

	defer func() {
		if r := recover(); r != nil {
			log.Println("Server exited!!!", r)
		}
	}()
}
