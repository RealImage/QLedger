package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	ledgerContext "github.com/RealImage/QLedger/context"
	"github.com/RealImage/QLedger/models"
)

func GetAccountsInfo(w http.ResponseWriter, r *http.Request, context *ledgerContext.AppContext) {
	accountsDB := models.Account{DB: context.DB}

	id := r.FormValue("id")
	account := accountsDB.GetByID(id)
	data, err := json.Marshal(account)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, string(data))
}
