package controllers

import (
	"bytes"
	"database/sql"
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
	TRANSACTIONS_API = "/v1/transactions"
)

type TransactionsSuite struct {
	suite.Suite
	context *ledgerContext.AppContext
}

func (ts *TransactionsSuite) SetupSuite() {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	assert.NotEmpty(ts.T(), databaseURL)
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Panic("Unable to connect to Database:", err)
	}
	log.Println("Successfully established connection to database.")
	ts.context = &ledgerContext.AppContext{DB: db}
}

func (ts *TransactionsSuite) TestValidTransaction() {
	t := ts.T()
	rr := httptest.NewRecorder()
	payload := `{
	  "id": "t001",
	  "lines": [
	    {
	      "account": "alice",
	      "delta": 100
	    },
	    {
	      "account": "bob",
	      "delta": -100
	    }
	  ]
	}`

	handler := middlewares.ContextMiddleware(MakeTransaction, ts.context)
	req, err := http.NewRequest("POST", TRANSACTIONS_API, bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "Invalid response code")
}

func (ts *TransactionsSuite) TestNoOpTransaction() {
	t := ts.T()
	rr := httptest.NewRecorder()
	payload := `{
	  "id": "t002",
	  "lines": []
	}`

	handler := middlewares.ContextMiddleware(MakeTransaction, ts.context)
	req, err := http.NewRequest("POST", TRANSACTIONS_API, bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "Invalid response code")
}

func (ts *TransactionsSuite) TestInvalidTransaction() {
	t := ts.T()
	rr := httptest.NewRecorder()
	payload := `{
	  "id": "t003",
	  "lines": [
	    {
	      "account": "alice",
	      "delta": 100
	    },
	    {
	      "account": "bob",
	      "delta": -101
	    }
	  ]
	}`

	handler := middlewares.ContextMiddleware(MakeTransaction, ts.context)
	req, err := http.NewRequest("POST", TRANSACTIONS_API, bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Invalid response code")
}

func (ts *TransactionsSuite) TestDuplicateTransaction() {
	t := ts.T()
	rr := httptest.NewRecorder()
	payload := `{
	  "id": "t001",
	  "lines": [
	    {
	      "account": "alice",
	      "delta": 100
	    },
	    {
	      "account": "bob",
	      "delta": -100
	    }
	  ]
	}`

	handler := middlewares.ContextMiddleware(MakeTransaction, ts.context)
	req, err := http.NewRequest("POST", TRANSACTIONS_API, bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code, "Invalid response code")
}

func (ts *TransactionsSuite) TestBadTransaction() {
	t := ts.T()
	rr := httptest.NewRecorder()
	payload := `{
		INVALID PAYLOAD
	}`

	handler := middlewares.ContextMiddleware(MakeTransaction, ts.context)
	req, err := http.NewRequest("POST", TRANSACTIONS_API, bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Invalid response code")
}

func (ts *TransactionsSuite) TearDownSuite() {
	// TODO: Cleanup test data
}

func TestTransactionsSuite(t *testing.T) {
	suite.Run(t, new(TransactionsSuite))
}
