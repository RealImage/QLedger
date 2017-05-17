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

func GetAccountsInfo(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	app := r.Context().Value("app").(*ledgerContext.AppContext)
	accountsDB := models.Account{DB: app.DB}

	accounts := accountsDB.GetByID()
	data, err := json.Marshal(accounts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, string(data))
}
