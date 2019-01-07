package client

import (
	"fmt"
	"testing"
)

func TestTransactions__Create(t *testing.T) {
	api := createTestAPI(t)
	api.AuthToken = testAuthToken

	// Grab account
	id := randomID()
	err := api.CreateTransaction(id, []Line{
		Line{Account: id + "1", Delta: -100},
		Line{Account: id + "2", Delta: 100},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Search for transactions
	transactions, err := api.GetTransactions(fmt.Sprintf(`{ "query": { "should": { "fields": [
            {"id": {"eq": "%s"}}
        ]}}}`, id))
	if err != nil {
		t.Fatal(err)
	}

	if len(transactions) != 1 {
		t.Errorf("got %d transactions", len(transactions))
	}

	tx := transactions[0]
	if tx.ID != id {
		t.Errorf("unknown transaction: %v", tx)
	}
	if len(tx.Lines) != 2 {
		t.Errorf("unknown Transaction.Lines: %v", tx.Lines)
	}
}
