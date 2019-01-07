package client

import (
	"fmt"
	"testing"
)

func TestAccounts__Create(t *testing.T) {
	api := createTestAPI(t)
	api.AuthToken = testAuthToken

	// Grab account
	id := randomID()
	err := api.CreateAccount(id, map[string]string{
		"test": "bar",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Search for account
	accounts, err := api.GetAccounts(fmt.Sprintf(`{ "query": { "should": { "fields": [
            {"id": {"eq": "%s"}}
        ]}}}`, id))
	if err != nil {
		t.Fatal(err)
	}

	if len(accounts) != 1 {
		t.Errorf("got %d accounts", len(accounts))
	}
	acct := accounts[0]
	if acct.ID != id {
		t.Errorf("unknown account: %v", acct)
	}
}
