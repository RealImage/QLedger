package middlewares

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/julienschmidt/httprouter"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthSuite struct {
	suite.Suite
	authServer *httptest.Server
}

func (as *AuthSuite) SetupSuite() {
	router := httprouter.New()
	router.HandlerFunc("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	})
	as.authServer = httptest.NewServer(TokenAuthMiddleware(router))
}

func (as *AuthSuite) TestNoAuth() {
	t := as.T()
	os.Unsetenv("LEDGER_AUTH_TOKEN")

	resp, err := http.Get(as.authServer.URL)
	if err != nil {
		log.Panic("Error connecting test server:", err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Invalid response code")
}

func (as *AuthSuite) TestValidAuth() {
	t := as.T()
	os.Setenv("LEDGER_AUTH_TOKEN", "XXX")

	req, _ := http.NewRequest("GET", as.authServer.URL, nil)
	req.Header.Add("LEDGER-AUTH-TOKEN", "XXX")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Panic("Error connecting test server:", err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Invalid response code")
}

func (as *AuthSuite) TestInvalidAuth() {
	t := as.T()
	os.Setenv("LEDGER_AUTH_TOKEN", "XXX")

	req, _ := http.NewRequest("GET", as.authServer.URL, nil)
	req.Header.Add("LEDGER-AUTH-TOKEN", "XYZ")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Panic("Error connecting test server:", err)
	}
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Invalid response code")
}

func (as *AuthSuite) TearDownSuite() {
	as.authServer.Close()
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}
