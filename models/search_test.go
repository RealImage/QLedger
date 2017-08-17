package models

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SearchSuite struct {
	suite.Suite
	db    *sql.DB
	accDB AccountDB
	txnDB TransactionDB
}

func (ss *SearchSuite) SetupSuite() {
	t := ss.T()
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	assert.NotEmpty(t, databaseURL)
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Panic("Unable to connect to Database:", err)
	} else {
		log.Println("Successfully established connection to database.")
		ss.db = db
	}
	log.Println("Successfully established connection to database.")
	ss.accDB = NewAccountDB(db)
	ss.txnDB = NewTransactionDB(db)

	// Create test accounts
	acc1 := &Account{
		ID: "acc1",
		Data: map[string]interface{}{
			"customer_id": "C1",
			"status":      "active",
			"created":     "2017-01-01",
		},
	}
	err = ss.accDB.CreateAccount(acc1)
	assert.Equal(t, nil, err, "Error creating test account")
	acc2 := &Account{
		ID: "acc2",
		Data: map[string]interface{}{
			"customer_id": "C2",
			"status":      "inactive",
			"created":     "2017-06-30",
		},
	}
	err = ss.accDB.CreateAccount(acc2)
	assert.Equal(t, nil, err, "Error creating test account")

	// Create test transactions
	txn1 := &Transaction{
		ID: "txn1",
		Lines: []*TransactionLine{
			&TransactionLine{
				AccountID: "acc1",
				Delta:     1000,
			},
			&TransactionLine{
				AccountID: "acc2",
				Delta:     -1000,
			},
		},
		Data: map[string]interface{}{
			"action": "setcredit",
			"expiry": "2018-01-01",
			"months": []string{"jan", "feb", "mar"},
		},
	}
	ok := ss.txnDB.Transact(txn1)
	assert.Equal(t, true, ok, "Error creating test transaction")
	txn2 := &Transaction{
		ID: "txn2",
		Lines: []*TransactionLine{
			&TransactionLine{
				AccountID: "acc1",
				Delta:     100,
			},
			&TransactionLine{
				AccountID: "acc2",
				Delta:     -100,
			},
		},
		Data: map[string]interface{}{
			"action": "setcredit",
			"expiry": "2018-01-15",
			"months": []string{"apr", "may", "jun"},
		},
	}
	ok = ss.txnDB.Transact(txn2)
	assert.Equal(t, true, ok, "Error creating test transaction")
	txn3 := &Transaction{
		ID: "txn3",
		Lines: []*TransactionLine{
			&TransactionLine{
				AccountID: "acc1",
				Delta:     400,
			},
			&TransactionLine{
				AccountID: "acc2",
				Delta:     -400,
			},
		},
		Data: map[string]interface{}{
			"action": "setcredit",
			"expiry": "2018-01-30",
			"months": []string{"jul", "aug", "sep"},
		},
	}
	ok = ss.txnDB.Transact(txn3)
	assert.Equal(t, true, ok, "Error creating test transaction")
}

func (ss *SearchSuite) TestSearchAccountsWithBothMustAndShould() {
	t := ss.T()
	engine, _ := NewSearchEngine(ss.db, "accounts")

	query := `{
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
	results, err := engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	accounts, _ := results.([]*AccountResult)
	assert.Equal(t, 1, len(accounts), "Account count doesn't match")
	assert.Equal(t, "acc1", accounts[0].ID, "Account ID doesn't match")
}

func (ss *SearchSuite) TestSearchTransactionsWithBothMustAndShould() {
	t := ss.T()
	engine, _ := NewSearchEngine(ss.db, "transactions")

	query := `{
        "query": {
            "must": {
                "fields": [
                    {"id": {"eq": "txn1"}}
                ],
                "terms": [
                    {"action": "setcredit"}
                ]
            },
            "should": {
                "terms": [
                    {"months": ["jan", "feb", "mar"]},
                    {"months": ["apr", "may", "jun"]},
                    {"months": ["jul", "aug", "sep"]}
                ],
                "ranges": [
                    {"expiry": {"gte": "2018-01-01", "lte": "2018-01-30"}}
                ]
            }
        }
    }`
	results, err := engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	transactions, _ := results.([]*TransactionResult)
	assert.Equal(t, 1, len(transactions), "Transaction count doesn't match")
	assert.Equal(t, "txn1", transactions[0].ID, "Transaction ID doesn't match")
}

func (ss *SearchSuite) TearDownSuite() {
	log.Println("Cleaning up the test database")

	t := ss.T()
	_, err := ss.db.Exec(`DELETE FROM lines`)
	if err != nil {
		t.Fatal("Error deleting lines:", err)
	}
	_, err = ss.db.Exec(`DELETE FROM transactions`)
	if err != nil {
		t.Fatal("Error deleting transactions:", err)
	}
	_, err = ss.db.Exec(`DELETE FROM accounts`)
	if err != nil {
		t.Fatal("Error deleting accounts:", err)
	}
}

func TestSearchSuite(t *testing.T) {
	suite.Run(t, new(SearchSuite))
}
