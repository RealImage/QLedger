package controllers

import (
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
	ACCOUNTS_INFO_API = "/v1/accounts"
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

func (as *AccountsSuite) TestAccountsInfoAPI() {
	t := as.T()
	rr := httptest.NewRecorder()

	handler := middlewares.ContextMiddleware(GetAccountInfo, as.context)
	req, err := http.NewRequest("GET", ACCOUNTS_INFO_API+"?id=100", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, 200, "Invalid response code")
	account := models.Account{}
	err = json.Unmarshal(rr.Body.Bytes(), &account)
	if err != nil {
		t.Errorf("Invalid json response: %v", rr.Body.String())
	}
	assert.Equal(t, account.Id, "100", "Invalid account ID")
	assert.Equal(t, account.Balance, 0, "Invalid account balance")
}

func TestAccountsSuite(t *testing.T) {
	suite.Run(t, new(AccountsSuite))
}
