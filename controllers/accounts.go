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

func GetAccounts(w http.ResponseWriter, r *http.Request, context *ledgerContext.AppContext) {
	defer r.Body.Close()
	engine, aerr := models.NewSearchEngine(context.DB, "accounts")
	if aerr != nil {
		log.Println("Error while creating Search Engine:", aerr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading payload:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := string(body)
	log.Println("Query:", query)

	results, aerr := engine.Query(query)
	if aerr != nil {
		log.Println("Error while querying:", aerr)
		switch aerr.ErrorCode() {
		case "search.query.invalid":
			w.WriteHeader(http.StatusBadRequest)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	data, err := json.Marshal(results)
	if err != nil {
		log.Println("Error while parsing results:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)

	return
}

func unmarshalToAccount(r *http.Request, account *models.Account) error {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, account)
	if err != nil {
		return err
	}
	var validKey = regexp.MustCompile(`^[a-z_A-Z]+$`)
	for key := range account.Data {
		if !validKey.MatchString(key) {
			return fmt.Errorf("Invalid key in data json: %v", key)
		}
	}
	return nil
}

func AddAccount(w http.ResponseWriter, r *http.Request, context *ledgerContext.AppContext) {
	account := &models.Account{}
	err := unmarshalToAccount(r, account)
	if err != nil {
		log.Println("Error loading payload:", err)
		w.WriteHeader(http.StatusBadRequest)
		//TODO Should we return any error message?
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
	account := &models.Account{}
	err := unmarshalToAccount(r, account)
	if err != nil {
		log.Println("Error loading payload:", err)
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
