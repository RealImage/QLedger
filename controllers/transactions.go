package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	ledgerContext "github.com/RealImage/QLedger/context"
	"github.com/RealImage/QLedger/models"
)

func unmarshalToTransaction(r *http.Request, txn *models.Transaction) error {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, txn)
	if err != nil {
		return err
	}
	var validKey = regexp.MustCompile(`^[a-z_A-Z]+$`)
	for key := range txn.Data {
		if !validKey.MatchString(key) {
			return fmt.Errorf("Invalid key in data json: %v", key)
		}
	}
	return nil
}

func MakeTransaction(w http.ResponseWriter, r *http.Request, context *ledgerContext.AppContext) {
	transaction := &models.Transaction{}
	err := unmarshalToTransaction(r, transaction)
	if err != nil {
		log.Println("Error loading payload:", err)
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

func GetTransactions(w http.ResponseWriter, r *http.Request, context *ledgerContext.AppContext) {
	defer r.Body.Close()
	engine, err := models.NewSearchEngine(context.DB, "transactions")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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
	if jerr != nil {
		log.Println("Error while parsing results:", jerr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(data))
	return
}

func UpdateTransaction(w http.ResponseWriter, r *http.Request, context *ledgerContext.AppContext) {
	transaction := &models.Transaction{}
	err := unmarshalToTransaction(r, transaction)
	if err != nil {
		log.Println("Error loading payload:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	transactionDB := models.TransactionDB{DB: context.DB}
	// Check if a transaction with same ID already exists
	if !transactionDB.IsExists(transaction.ID) {
		log.Println("Transaction doesn't exist:", transaction.ID)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Otherwise, update transaction
	terr := transactionDB.UpdateTransaction(transaction)
	if terr == nil {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		log.Printf("Error while updating transaction: %v (%v)", transaction.ID, terr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
