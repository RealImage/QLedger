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
	TransactionsAPI = "/v1/transactions"
)

type TransactionsSuite struct {
	suite.Suite
	context *ledgerContext.AppContext
}

func (ts *TransactionsSuite) SetupSuite() {
	log.Println("Connecting to the test database")
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	assert.NotEmpty(ts.T(), databaseURL)
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Panic("Unable to connect to Database:", err)
	}
	log.Println("Successfully established connection to database.")
	ts.context = &ledgerContext.AppContext{DB: db}
}

func (ts *TransactionsSuite) TestValidAndRepeatedTransaction() {
	t := ts.T()

	// Valid transaction
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
	  ],
	  "data": {
	    "tag_one": "val1",
	    "tag_two": "val2"
	  }
	}`
	handler := middlewares.ContextMiddleware(MakeTransaction, ts.context)
	req, err := http.NewRequest("POST", TransactionsAPI, bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	rr1 := httptest.NewRecorder()
	handler.ServeHTTP(rr1, req)
	assert.Equal(t, http.StatusCreated, rr1.Code, "Invalid response code")

	// Duplicate transaction
	req, err = http.NewRequest("POST", TransactionsAPI, bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req)
	assert.Equal(t, http.StatusAccepted, rr2.Code, "Invalid response code")

	// Conflict transaction
	payload = `{
	  "id": "t001",
	  "lines": [
	    {
	      "account": "alice",
	      "delta": 200
	    },
	    {
	      "account": "bob",
	      "delta": -200
	    }
	  ]
	}`
	req, err = http.NewRequest("POST", TransactionsAPI, bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	rr3 := httptest.NewRecorder()
	handler.ServeHTTP(rr3, req)
	assert.Equal(t, http.StatusConflict, rr3.Code, "Invalid response code")
}

func (ts *TransactionsSuite) TestNoOpTransaction() {
	t := ts.T()
	rr := httptest.NewRecorder()
	payload := `{
	  "id": "t002",
	  "lines": []
	}`

	handler := middlewares.ContextMiddleware(MakeTransaction, ts.context)
	req, err := http.NewRequest("POST", TransactionsAPI, bytes.NewBufferString(payload))
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
	req, err := http.NewRequest("POST", TransactionsAPI, bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Invalid response code")
}

func (ts *TransactionsSuite) TestBadTransaction() {
	t := ts.T()
	rr := httptest.NewRecorder()
	payload := `{
		INVALID PAYLOAD
	}`

	handler := middlewares.ContextMiddleware(MakeTransaction, ts.context)
	req, err := http.NewRequest("POST", TransactionsAPI, bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Invalid response code")
}

func (ts *TransactionsSuite) TestFailTransaction() {
	t := ts.T()
	rr := httptest.NewRecorder()
	payload := `{
	  "id": "t004",
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

	// database is not available
	db, _ := sql.Open("postgres", "")
	invalidContext := &ledgerContext.AppContext{DB: db}

	handler := middlewares.ContextMiddleware(MakeTransaction, invalidContext)
	req, err := http.NewRequest("POST", TransactionsAPI, bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Invalid response code")
}

func (ts *TransactionsSuite) TearDownSuite() {
	log.Println("Cleaning up the test database")

	t := ts.T()
	_, err := ts.context.DB.Exec(`DELETE FROM lines`)
	if err != nil {
		t.Fatal("Error deleting lines:", err)
	}
	_, err = ts.context.DB.Exec(`DELETE FROM transactions`)
	if err != nil {
		t.Fatal("Error deleting transactions:", err)
	}
	_, err = ts.context.DB.Exec(`DELETE FROM accounts`)
	if err != nil {
		t.Fatal("Error deleting accounts:", err)
	}
}

func TestTransactionsSuite(t *testing.T) {
	suite.Run(t, new(TransactionsSuite))
}
