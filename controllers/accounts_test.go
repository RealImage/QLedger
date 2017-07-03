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
	ACCOUNTS_API = "/v1/accounts"
)

type AccountsSuite struct {
	suite.Suite
	context *ledgerContext.AppContext
}

func (as *AccountsSuite) SetupTest() {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	assert.NotEmpty(as.T(), databaseURL)
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Panic("Unable to connect to Database:", err)
	}
	log.Println("Successfully established connection to database.")
	as.context = &ledgerContext.AppContext{DB: db}
}

func (as *AccountsSuite) TestAccountsAPI() {
	t := as.T()

	// Sample account which doesn't exist
	payload := `{
	  "query": {
	    "id": "acc1"
	  }
	}`
	handler := middlewares.ContextMiddleware(GetAccounts, as.context)
	req, err := http.NewRequest("GET", ACCOUNTS_API, bytes.NewBufferString(payload))
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
	assert.Equal(t, len(accounts), 0, "Account should not exist")
}

func TestAccountsSuite(t *testing.T) {
	suite.Run(t, new(AccountsSuite))
}
