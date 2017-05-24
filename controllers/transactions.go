package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	ledgerContext "github.com/RealImage/QLedger/context"
	"github.com/RealImage/QLedger/models"
)

func MakeTransaction(w http.ResponseWriter, r *http.Request, context *ledgerContext.AppContext) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading payload:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	transaction := &models.Transaction{}
	err = json.Unmarshal(body, transaction)
	if err != nil {
		log.Println("Error loading JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	transactionsDB := models.TransactionDB{DB: context.DB}
	// Skip if transaction already exists
	if transactionsDB.IsExists(transaction.ID) {
		log.Println("Transaction is duplicate:", transaction.ID)
		w.WriteHeader(http.StatusConflict)
		return
	}
	// Skip if the transaction is invalid
	// by validating the delta values
	if !transaction.IsValid() {
		log.Println("Transaction is invalid:", transaction.ID)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Otherwise, do transaction
	done := transactionsDB.DoTransaction(transaction)
	if done {
		w.WriteHeader(http.StatusCreated)
		return
	} else {
		log.Println("Transaction failed:", transaction.ID)
		w.WriteHeader(http.StatusAccepted)
		return
	}
}
