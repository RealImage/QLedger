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

type MockAccount struct {
	Id      string `json:"id"`
	Balance int    `json:"balance"`
}

func (as *AccountsSuite) SetupTest() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Panic("Unable to connect to Database:", err)
	}
	log.Println("Successfully established connection to database.")
	as.context = &ledgerContext.AppContext{DB: db}
}

func (as *AccountsSuite) TestValidAccount() {
	t := as.T()
	rr := httptest.NewRecorder()

	handler := middlewares.ContextMiddleware(GetAccountInfo, as.context)
	req, err := http.NewRequest("GET", ACCOUNTS_INFO_API+"?id=100", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, req)

	account := MockAccount{}
	// test valid status code
	assert.Equal(t, rr.Code, 200, "Invalid response code")
	// test valid json
	err = json.Unmarshal(rr.Body.Bytes(), &account)
	if err != nil {
		t.Errorf("Invalid json response: %v", rr.Body.String())
	}
	// test valid id
	assert.Equal(t, account.Id, "100", "Invalid account ID")
	// test valid balance
	assert.Equal(t, account.Balance, 5, "Invalid account balance")
}

func (as *AccountsSuite) TestInvalidAccount() {
	t := as.T()
	rr := httptest.NewRecorder()

	handler := middlewares.ContextMiddleware(GetAccountInfo, as.context)
	req, err := http.NewRequest("GET", ACCOUNTS_INFO_API+"?id=101", nil)
	if err != nil {
		as.T().Fatal(err)
	}
	handler.ServeHTTP(rr, req)

	account := MockAccount{}
	// test valid status code
	assert.Equal(t, rr.Code, 200, "Invalid response code")
	// test valid json
	err = json.Unmarshal(rr.Body.Bytes(), &account)
	if err != nil {
		t.Errorf("Invalid json response: %v", rr.Body.String())
	}
	// test valid id
	assert.Equal(t, account.Id, "101", "Invalid account ID")
	// test valid balance
	assert.Equal(t, account.Balance, 0, "Invalid account balance")
}

func TestAccountsSuite(t *testing.T) {
	suite.Run(t, new(AccountsSuite))
}
