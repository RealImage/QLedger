package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (a *API) CreateAccount(id string, data map[string]string) error {
	wrapper := make(map[string]interface{}, 0)
	wrapper["id"] = id
	wrapper["data"] = data

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(wrapper); err != nil {
		return fmt.Errorf("problem encoding %s Account JSON: %v", id, err)
	}

	req, err := http.NewRequest("POST", a.Host+a.buildPath("/accounts"), &buf)
	if err != nil {
		return fmt.Errorf("unable to build CreateAccount request: %v", err)
	}
	a.setAuthToken(req)

	resp, err := a.Underlying.Do(req)
	if err != nil {
		return fmt.Errorf("problem making CreateAccount: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return fmt.Errorf("bogus error code in CreateAccount: %s", resp.Status)
	}

	return nil
}

type Account struct {
	ID      string            `json:"id"`
	Balance int64             `json:"balance"`
	Data    map[string]string `json:"data"`
}

func (a Account) String() string {
	return fmt.Sprintf("Account=%q Balance=%q Metadata=%#v", a.ID, a.Balance, a.Data)
}

func (a *API) GetAccounts(body string) ([]Account, error) {
	req, err := http.NewRequest("GET", a.Host+a.buildPath("/accounts"), strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("unable to build GetAccounts request: %v", err)
	}
	a.setAuthToken(req)

	resp, err := a.Underlying.Do(req)
	if err != nil {
		return nil, fmt.Errorf("problem making GetAccounts: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("bogus error code in GetAccounts: %s", resp.Status)
	}

	// Read response body
	var accounts []Account
	if err := json.NewDecoder(resp.Body).Decode(&accounts); err != nil {
		return nil, fmt.Errorf("unable to read GetAccounts response body: %v", err)
	}
	return accounts, nil
}
