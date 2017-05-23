package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	ledgerContext "github.com/RealImage/QLedger/context"
	"github.com/RealImage/QLedger/controllers"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

type Handler func(http.ResponseWriter, *http.Request, *ledgerContext.AppContext)

func ContextMiddleware(handler Handler, context *ledgerContext.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, context)
	}
}

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Panic("Unable to connect to Database:", err)
	}
	log.Println("Successfully established connection to database.")
	appContext := &ledgerContext.AppContext{DB: db}

	router := httprouter.New()
	router.HandlerFunc("GET", "/v1/accounts", ContextMiddleware(controllers.GetAccountsInfo, appContext))
	router.HandlerFunc("POST", "/v1/transactions", ContextMiddleware(controllers.MakeTransaction, appContext))

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
