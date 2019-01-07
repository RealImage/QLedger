package client

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"
	"strings"
	"testing"
)

var (
	testHost, testPath = os.Getenv("TEST_HOST"), os.Getenv("TEST_PATH")
	testAuthToken      = os.Getenv("TEST_AUTH_TOKEN")
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

func randomID() string {
	bs := make([]byte, 20)
	n, err := rand.Read(bs)
	if err != nil || n == 0 {
		return ""
	}
	return strings.ToLower(hex.EncodeToString(bs))
}

func TestAPI__Ping(t *testing.T) {
	api := createTestAPI(t)
	if err := api.Ping(); err != nil {
		t.Error(err)
	}
}
