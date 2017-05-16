package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/RealImage/QLedger/models"
	"github.com/julienschmidt/httprouter"
)

func GetAccountsInfo(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	accounts := models.GetAccounts()
	data, err := json.Marshal(accounts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, string(data))
}
