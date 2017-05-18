package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	ledgerContext "github.com/RealImage/QLedger/context"
	"github.com/RealImage/QLedger/models"
	"github.com/julienschmidt/httprouter"
)

func GetAccountsInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	app := r.Context().Value("app").(*ledgerContext.AppContext)
	accountsDB := models.Account{DB: app.DB}
	
	id := r.FormValue("id")
	account := accountsDB.GetByID(id)
	data, err := json.Marshal(account)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, string(data))
}
