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

func AddAccount(w http.ResponseWriter, r *http.Request, context *ledgerContext.AppContext) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading payload:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	account := &models.Account{}
	err = json.Unmarshal(body, account)
	if err != nil {
		log.Println("Error loading JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	accountsDB := models.AccountDB{DB: context.DB}
	// Check if an account with same ID already exists
	if accountsDB.IsExists(account.Id) {
		log.Println("Account is conflicting:", account.Id)
		w.WriteHeader(http.StatusConflict)
		return
	}

	// Otherwise, add account
	aerr := accountsDB.CreateAccount(account)
	if aerr == nil {
		w.WriteHeader(http.StatusCreated)
		return
	} else {
		log.Printf("Error while adding account: %v (%v)", account.Id, aerr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func UpdateAccount(w http.ResponseWriter, r *http.Request, context *ledgerContext.AppContext) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading payload:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	account := &models.Account{}
	err = json.Unmarshal(body, account)
	if err != nil {
		log.Println("Error loading JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	accountsDB := models.AccountDB{DB: context.DB}
	// Check if an account with same ID already exists
	if !accountsDB.IsExists(account.Id) {
		log.Println("Account doesn't exist:", account.Id)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Otherwise, update account
	aerr := accountsDB.UpdateAccount(account)
	if aerr == nil {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		log.Printf("Error while updating account: %v (%v)", account.Id, aerr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
