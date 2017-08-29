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

	// Migrate DB changes
	migrateDB(db)

	appContext := &ledgerContext.AppContext{DB: db}
	router := httprouter.New()

	// Monitors
	router.HandlerFunc(http.MethodGet, "/ping", controllers.Ping)

	// Create accounts and transactions
	router.HandlerFunc(http.MethodPost, "/v1/accounts", middlewares.ContextMiddleware(controllers.AddAccount, appContext))
	router.HandlerFunc(http.MethodPost, "/v1/transactions", middlewares.ContextMiddleware(controllers.MakeTransaction, appContext))

	// Read or search accounts and transactions
	router.HandlerFunc(http.MethodGet, "/v1/accounts", middlewares.ContextMiddleware(controllers.GetAccounts, appContext))
	router.HandlerFunc(http.MethodPost, "/v1/accounts/_search", middlewares.ContextMiddleware(controllers.GetAccounts, appContext))
	router.HandlerFunc(http.MethodGet, "/v1/transactions", middlewares.ContextMiddleware(controllers.GetTransactions, appContext))
	router.HandlerFunc(http.MethodPost, "/v1/transactions/_search", middlewares.ContextMiddleware(controllers.GetTransactions, appContext))

	// Update data of accounts and transactions
	router.HandlerFunc(http.MethodPut, "/v1/accounts", middlewares.ContextMiddleware(controllers.UpdateAccount, appContext))
	router.HandlerFunc(http.MethodPut, "/v1/transactions", middlewares.ContextMiddleware(controllers.UpdateTransaction, appContext))

	port := os.Getenv("PORT")
	if port == "" {
		port = "7000"
	}
	log.Println("Running server on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, middlewares.TokenAuthMiddleware(router)))

	defer func() {
		if r := recover(); r != nil {
			log.Println("Server exited!!!", r)
		}
	}()
}

func migrateDB(db *sql.DB) {
	log.Println("Starting db schema migration...")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Panic("Unable to create database instance for migration:", err)
	}

	migrationFilesPath := os.Getenv("MIGRATION_FILES_PATH")
	if migrationFilesPath == "" {
		migrationFilesPath = "file://migrations/postgres"
	}
	m, err := migrate.NewWithDatabaseInstance(
		migrationFilesPath,
		"postgres", driver)
	if err != nil {
		log.Panic("Unable to create Migrate instance for database:", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Panic("Unable to get existing migration version for database:", dirty, err)
	}
	log.Println("Current schema version:", version)
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange && err != migrate.ErrLocked {
		log.Println("Error while migration:", err)
	}
	version, dirty, err = m.Version()
	if err != nil {
		log.Panic("Unable to get new migration version for database:", dirty, err)
	}
	log.Println("Migrated schema version:", version)
}
