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
	"github.com/mattes/migrate/database"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

func main() {
	// Assert authentication
	authToken, ok := os.LookupEnv("LEDGER_AUTH_TOKEN")
	if !ok || authToken == "" {
		log.Fatal("Cannot start the server. Authentication token is not set!! Please set LEDGER_AUTH_TOKEN")
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Panic("Unable to connect to Database:", err)
	}
	log.Println("Successfully established connection to database.")

	// Migrate DB changes
	migrateDB(db)

	appContext := &ledgerContext.AppContext{DB: db}
	router := httprouter.New()

	hostPrefix := os.Getenv("HOST_PREFIX")
	// Monitors
	router.HandlerFunc(http.MethodGet, hostPrefix+"/ping", controllers.Ping)

	// Create accounts and transactions
	router.HandlerFunc(http.MethodPost, hostPrefix+"/v1/accounts",
		middlewares.TokenAuthMiddleware(
			middlewares.ContextMiddleware(controllers.AddAccount, appContext)))
	router.HandlerFunc(http.MethodPost, hostPrefix+"/v1/transactions",
		middlewares.TokenAuthMiddleware(
			middlewares.ContextMiddleware(controllers.MakeTransaction, appContext)))

	// Read or search accounts and transactions
	router.HandlerFunc(http.MethodGet, hostPrefix+"/v1/accounts",
		middlewares.TokenAuthMiddleware(
			middlewares.ContextMiddleware(controllers.GetAccounts, appContext)))
	router.HandlerFunc(http.MethodPost, hostPrefix+"/v1/accounts/_search",
		middlewares.TokenAuthMiddleware(
			middlewares.ContextMiddleware(controllers.GetAccounts, appContext)))
	router.HandlerFunc(http.MethodGet, hostPrefix+"/v1/transactions",
		middlewares.TokenAuthMiddleware(
			middlewares.ContextMiddleware(controllers.GetTransactions, appContext)))
	router.HandlerFunc(http.MethodPost, hostPrefix+"/v1/transactions/_search",
		middlewares.TokenAuthMiddleware(
			middlewares.ContextMiddleware(controllers.GetTransactions, appContext)))

	// Update data of accounts and transactions
	router.HandlerFunc(http.MethodPut, hostPrefix+"/v1/accounts",
		middlewares.TokenAuthMiddleware(
			middlewares.ContextMiddleware(controllers.UpdateAccount, appContext)))
	router.HandlerFunc(http.MethodPut, hostPrefix+"/v1/transactions",
		middlewares.TokenAuthMiddleware(
			middlewares.ContextMiddleware(controllers.UpdateTransaction, appContext)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "7000"
	}
	log.Println("Running server on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, router))

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
	if err != nil {
		switch err {
		case migrate.ErrNoChange:
			log.Println("No changes to migrate")
		case migrate.ErrLocked, database.ErrLocked:
			log.Println("Database locked. Skipping migration assuming another instance working on it")
		default:
			log.Panic("Error while migration:", err)
		}
	}
	version, dirty, err = m.Version()
	if err != nil {
		log.Panic("Unable to get new migration version for database:", dirty, err)
	}
	log.Println("Migrated schema version:", version)
}
