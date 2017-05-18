package main

import (
	"log"
	"net/http"
	"os"

	ledgerContext "github.com/RealImage/QLedger/context"
	"github.com/RealImage/QLedger/controllers"
	"github.com/julienschmidt/httprouter"
)

type Handler func(http.ResponseWriter, *http.Request, *ledgerContext.AppContext)

func ContextMiddleware(handler Handler, context *ledgerContext.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, context)
	}
}

func main() {
	appContext := &ledgerContext.AppContext{}
	appContext.Initialize()

	router := httprouter.New()
	router.HandlerFunc("GET", "/v1/accounts", ContextMiddleware(controllers.GetAccountsInfo, appContext))

	port := "7000"
	if lp := os.Getenv("PORT"); lp != "" {
		port = lp
	}
	log.Println("Running server on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, router))

	defer func() {
		appContext.Cleanup()
		if r := recover(); r != nil {
			log.Println("Server exited!!!", r)
		}
	}()
}
