package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	ledgerContext "github.com/RealImage/QLedger/context"
	"github.com/RealImage/QLedger/models"
	"github.com/julienschmidt/httprouter"
)

func Search(w http.ResponseWriter, r *http.Request, params httprouter.Params, context *ledgerContext.AppContext) {
	namespace := params.ByName("namespace")

	engine, err := models.NewSearchEngine(namespace)
	if err != nil {
		log.Println("Error initiating search:", err)
		switch err.ErrorCode() {
		case "search.namespace.invalid":
			w.WriteHeader(http.StatusNotFound)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	defer r.Body.Close()
	body, rerr := ioutil.ReadAll(r.Body)
	if rerr != nil {
		log.Println("Error reading payload:", rerr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	query := string(body)

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
	w.WriteHeader(http.StatusOK)
}
