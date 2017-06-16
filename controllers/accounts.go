package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	ledgerContext "github.com/RealImage/QLedger/context"
	"github.com/RealImage/QLedger/models"
)

func GetAccounts(w http.ResponseWriter, r *http.Request, context *ledgerContext.AppContext) {
	engine, err := models.NewSearchEngine(context.DB, "accounts")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()
	body, rerr := ioutil.ReadAll(r.Body)
	if rerr != nil {
		log.Println("Error reading payload:", rerr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	query := string(body)
	log.Println("Query:", query)

	results, err := engine.Query(query)
	if err != nil {
		log.Println("Error while querying:", err)
		switch err.ErrorCode() {
		case "search.query.invalid":
			w.WriteHeader(http.StatusBadRequest)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	data, jerr := json.Marshal(results)
	if err != nil {
		log.Println("Error while parsing results:", jerr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(data))
	return
}
