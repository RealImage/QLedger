package main

import (
	"log"
	"net/http"

	"github.com/RealImage/QLedger/controllers"
	"github.com/RealImage/QLedger/database"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/v1/accounts", controllers.GetAccountsInfo)

	port := "7000" //TODO: Read from config
	log.Println("Running server on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, router))

	defer func() {
		database.Cleanup()
		if r := recover(); r != nil {
			log.Println("Server exited!!!", r)
		}
	}()
}
