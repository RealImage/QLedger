package main

import (
	"context"
	"log"
	"net/http"

	ledgerContext "github.com/RealImage/QLedger/context"
	"github.com/RealImage/QLedger/controllers"
	"github.com/julienschmidt/httprouter"
)

func ContextMiddleware(appContext *ledgerContext.AppContext, handlerFunc httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), "app", appContext)
		handlerFunc(w, r.WithContext(ctx), ps)
	}
}

func main() {
	appContext := &ledgerContext.AppContext{}
	appContext.Initialize()

	router := httprouter.New()
	router.GET("/v1/accounts", ContextMiddleware(appContext, controllers.GetAccountsInfo))

	port := "7000" //TODO: Read from config
	log.Println("Running server on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, router))

	defer func() {
		appContext.Cleanup()
		if r := recover(); r != nil {
			log.Println("Server exited!!!", r)
		}
	}()
}
