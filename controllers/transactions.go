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
	// Check if a transaction with same ID already exists
	if transactionsDB.IsExists(transaction.ID) {
		// Check if the transaction lines are different
		// and conflicts with the existing lines
		if transactionsDB.IsConflict(transaction) {
			// The conflicting transactions are denied
			log.Println("Transaction is conflicting:", transaction.ID)
			w.WriteHeader(http.StatusConflict)
			return
		} else {
			// Otherwise the transaction is just a duplicate
			// The exactly duplicate transactions are ignored
			log.Println("Transaction is duplicate:", transaction.ID)
			w.WriteHeader(http.StatusAccepted)
			return
		}
	}

	// Skip if the transaction is invalid
	// by validating the delta values
	if !transaction.IsValid() {
		log.Println("Transaction is invalid:", transaction.ID)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Otherwise, do transaction
	done := transactionsDB.Transact(transaction)
	if done {
		w.WriteHeader(http.StatusCreated)
		return
	} else {
		log.Println("Transaction failed:", transaction.ID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
