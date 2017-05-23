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

func MakeTransaction(w http.ResponseWriter, r *http.Request, context *ledgerContext.AppContext) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading payload:", err)
		return
	}
	transaction := &models.Transaction{}
	err = json.Unmarshal(body, transaction)
	if err != nil {
		log.Println("Error loading JSON:", err)
		return
	}

	transactionsDB := models.TransactionDB{DB: context.DB}
	// Skip if transaction already exists
	if transactionsDB.IsExists(transaction.ID) {
		fmt.Fprint(w, "Transaction already exists")
		return
	}
	// Skip if the transaction is invalid
	// by validating the delta values
	if !transaction.IsValid() {
		fmt.Fprint(w, "Transaction invalid")
		return
	}

	// Otherwise, do transaction
	done := transactionsDB.DoTransaction(transaction)
	if done {
		fmt.Fprint(w, "Transaction is success")
		return
	} else {
		fmt.Fprint(w, "Transaction failed")
		return
	}

	//TODO: Return JSON response with proper status codes
}
