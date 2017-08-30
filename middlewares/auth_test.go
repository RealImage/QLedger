package middlewares

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthSuite struct {
	suite.Suite
	handler http.HandlerFunc
}

func (as *AuthSuite) SetupSuite() {
	as.handler = TokenAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	})
}

func (as *AuthSuite) TestNoAuth() {
	t := as.T()
	os.Unsetenv("LEDGER_AUTH_TOKEN")

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr1 := httptest.NewRecorder()
	as.handler.ServeHTTP(rr1, req)
	assert.Equal(t, http.StatusOK, rr1.Code, "Invalid response code")
}

func (as *AuthSuite) TestValidAuth() {
	t := as.T()
	os.Setenv("LEDGER_AUTH_TOKEN", "XXX")

	req, err := http.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "XXX")
	if err != nil {
		t.Fatal(err)
	}
	rr1 := httptest.NewRecorder()
	as.handler.ServeHTTP(rr1, req)
	assert.Equal(t, http.StatusOK, rr1.Code, "Invalid response code")
}

func (as *AuthSuite) TestInvalidAuth() {
	t := as.T()
	os.Setenv("LEDGER_AUTH_TOKEN", "XXX")

	req, err := http.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "XYZ")
	if err != nil {
		t.Fatal(err)
	}
	rr1 := httptest.NewRecorder()
	as.handler.ServeHTTP(rr1, req)
	assert.Equal(t, http.StatusUnauthorized, rr1.Code, "Invalid response code")
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}
