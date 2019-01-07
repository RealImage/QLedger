package client

import (
	"net/http"
	"os"
	"testing"
)

var (
	testHost, testPath = os.Getenv("TEST_HOST"), os.Getenv("TEST_PATH")
)

func TestAPI__BasePath(t *testing.T) {
	api, err := New("sbx.example.com", http.DefaultClient)
	if err != nil {
		t.Fatal(err)
	}

	// Defaults
	if v := api.buildPath("/transactions"); v != "/v1/transactions" {
		t.Errorf("got %q", v)
	}

	// Example of hiding QLedger behind a LB
	api.BasePath = "/v1/ledger"
	if v := api.buildPath("/transactions"); v != "/v1/ledger/transactions" {
		t.Errorf("got %q", v)
	}
}

func createTestAPI(t *testing.T) *API {
	t.Helper()

	if testHost == "" {
		t.Skip("missing TEST_HOST to make network calls")
	}

	api, err := New(testHost, http.DefaultClient)
	if err != nil {
		t.Fatal(err)
	}
	if testPath != "" {
		api.BasePath = testPath
	}

	return api
}

func TestAPI__Ping(t *testing.T) {
	api := createTestAPI(t)
	if err := api.Ping(); err != nil {
		t.Error(err)
	}
}
