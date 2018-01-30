package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	ledgerContext "github.com/RealImage/QLedger/context"
	"github.com/RealImage/QLedger/models"
)

const (
	// TimestampQueryLayout is the timestamp layout in balance queries
	TimestampQueryLayout = "2006-01-02 15:04:05"
)

// BalanceRequest ...
type BalanceRequest struct {
	AccountID string `json:"account_id"`
	On        string `json:"on"`
}

// BalanceResponse ...
type BalanceResponse struct {
	AccountID string `json:"account_id"`
	On        string `json:"on"`
	Balance   int    `json:"balance"`
}

// GetBalanceOnTime returns the balance of an account on a particular time
func GetBalanceOnTime(w http.ResponseWriter, r *http.Request, context *ledgerContext.AppContext) {
	// Parse the input data
	defer r.Body.Close()
	inputData := &BalanceRequest{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, inputData)
	if err != nil {
		log.Println("Error loading JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = time.Parse(TimestampQueryLayout, inputData.On)
	if err != nil {
		log.Println("Error while validating timestamp:", inputData.On, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the balance as on time from database
	accountsDB := models.NewAccountDB(context.DB)
	acc, err := accountsDB.GetBalanceOnTime(inputData.AccountID, inputData.On)
	if err != nil {
		log.Printf("Error while getting balance of account: %v as on time: %v (%v)",
			inputData.AccountID, inputData.On, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Prepare the response
	response := &BalanceResponse{
		AccountID: inputData.AccountID,
		On:        inputData.On,
		Balance:   acc.Balance,
	}
	data, err := json.Marshal(response)
	if err != nil {
		log.Println("Error while building response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
	return
}
