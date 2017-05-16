package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RealImage/QLedger/internal/utils/config"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

func AddTransaction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Add Transaction - WIP!\n")
}

func GetAccountInformation(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Get Account Information - WIP!\n")
}

func main() {
	router := httprouter.New()
	router.POST("/v1/transactions", AddTransaction)
	router.GET("/v1/accounts", GetAccountInformation)

	log.Println("PORT: ", config.PORT)
	log.Fatal(http.ListenAndServe(":"+config.PORT, context.ClearHandler(router)))
}
