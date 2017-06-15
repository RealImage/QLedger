package controllers

import (
	"io/ioutil"
	"log"
	"net/http"

	ledgerContext "github.com/RealImage/QLedger/context"
	"github.com/julienschmidt/httprouter"
)

func Search(w http.ResponseWriter, r *http.Request, params httprouter.Params, context *ledgerContext.AppContext) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading payload:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	query := string(body)
	namespace := params.ByName("namespace")

	log.Println("query:", query)
	log.Println("namespace:", namespace)
	w.WriteHeader(http.StatusOK)
}
