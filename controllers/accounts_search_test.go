package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	ledgerContext "github.com/RealImage/QLedger/context"
	"github.com/RealImage/QLedger/middlewares"
	"github.com/RealImage/QLedger/models"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	ACCOUNTS_SEARCH_API = "/v1/accounts"
)

type AccountsSearchSuite struct {
	suite.Suite
	context *ledgerContext.AppContext
}

func (as *AccountsSearchSuite) SetupTest() {
	t := as.T()
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	assert.NotEmpty(t, databaseURL)
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Panic("Unable to connect to Database:", err)
	}
	log.Println("Successfully established connection to database.")
	as.context = &ledgerContext.AppContext{DB: db}

	// Create test accounts
	accDB := models.AccountDB{DB: db}
	acc1 := &models.Account{
		Id: "acc1",
		Data: map[string]interface{}{
			"customer_id": "C1",
			"status":      "active",
			"created":     "2017-01-01",
		},
	}
	err = accDB.CreateAccount(acc1)
	assert.Equal(t, nil, err, "Error creating test account")
	acc2 := &models.Account{
		Id: "acc2",
		Data: map[string]interface{}{
			"customer_id": "C2",
			"status":      "inactive",
			"created":     "2017-06-30",
		},
	}
	err = accDB.CreateAccount(acc2)
	assert.Equal(t, nil, err, "Error creating test account")
}

func (as *AccountsSearchSuite) TestAccountsSearch() {
	t := as.T()

	// Prepare search query
	payload := `{
        "query": {
            "must": {
                "fields": [
                    {"id": {"eq": "acc1"}}
                ],
                "terms": [
                    {"status": "active"}
                ]
            },
            "should": {
                "terms": [
                    {"customer_id": "C1"}
                ],
                "ranges": [
                    {"created": {"gte": "2018-01-01", "lte": "2018-01-30"}}
                ]
            }
        }
    }`
	handler := middlewares.ContextMiddleware(GetAccounts, as.context)
	req, err := http.NewRequest("GET", ACCOUNTS_SEARCH_API, bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "Invalid response code")

	var accounts []models.AccountResult
	err = json.Unmarshal(rr.Body.Bytes(), &accounts)
	if err != nil {
		t.Errorf("Invalid json response: %v", rr.Body.String())
	}
	assert.Equal(t, 1, len(accounts), "Accounts count doesn't match")
	assert.Equal(t, "acc1", accounts[0].ID, "Account ID doesn't match")
}

func TestAccountsSuite(t *testing.T) {
	suite.Run(t, new(AccountsSearchSuite))
}
