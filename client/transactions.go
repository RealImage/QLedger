package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Line struct {
	Account string
	Delta   int64
}

func (l Line) String() string {
	return fmt.Sprintf("Line: Account=%q Delta=%d", l.Account, l.Delta)
}

func (a *API) CreateTransaction(id string, lines []Line) error {
	wrapper := make(map[string]interface{}, 0)
	wrapper["id"] = id
	wrapper["lines"] = lines

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(wrapper); err != nil {
		return fmt.Errorf("problem encoding %s Transaction JSON: %v", id, err)
	}

	req, err := http.NewRequest("POST", a.Host+a.buildPath("/transactions"), &buf)
	if err != nil {
		return fmt.Errorf("unable to build CreateTransaction request: %v", err)
	}
	a.setAuthToken(req)

	resp, err := a.Underlying.Do(req)
	if err != nil {
		return fmt.Errorf("problem making CreateTransaction: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return fmt.Errorf("bogus error code in CreateTransaction: %s", resp.Status)
	}

	return nil
}

type Transaction struct {
	ID        string            `json:"id"`
	Timestamp time.Time         `json:"timestamp"`
	Data      map[string]string `json:"data"`
	Lines     []Line            `json:"lines"`
}

func (t Transaction) String() string {
	return fmt.Sprintf("Transaction=%q Timestamp=%q Data=%q Lines=%v", t.ID, t.Timestamp, t.Data, t.Lines)
}

func (a *API) GetTransactions(body string) ([]Transaction, error) {
	req, err := http.NewRequest("GET", a.Host+a.buildPath("/transactions"), strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("unable to build GetTransactions request: %v", err)
	}
	a.setAuthToken(req)

	resp, err := a.Underlying.Do(req)
	if err != nil {
		return nil, fmt.Errorf("problem making GetTransactions: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("bogus error code in GetTransactions: %s", resp.Status)
	}

	// Read response body
	var transactions []Transaction
	if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
		return nil, fmt.Errorf("unable to read GetTransactions response body: %v", err)
	}
	return transactions, nil
}
