package client

import (
	"fmt"
	"net/http"
	"path"
	"strings"
)

func New(host string, underlying *http.Client) (*API, error) {
	if host == "" && !strings.HasPrefix("http", host) {
		return nil, fmt.Errorf("invalid host: %q", host)
	}

	return &API{
		BasePath:   "/v1/",
		Host:       host,
		Underlying: underlying,
	}, nil
}

type API struct {
	Host, BasePath string
	AuthToken      string

	Underlying *http.Client
}

func (a *API) buildPath(sub string) string {
	return path.Join(a.BasePath, sub)
}

func (a *API) Ping() error {
	req, err := http.NewRequest("GET", a.Host+a.buildPath("/ping"), nil)
	if err != nil {
		return fmt.Errorf("unable to build Ping: %v", err)
	}
	resp, err := a.Underlying.Do(req)
	if err != nil {
		return fmt.Errorf("problem making Ping: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return fmt.Errorf("bogus error code: %s", resp.Status)
	}

	return nil
}

func (a *API) setAuthToken(req *http.Request) {
	if a.AuthToken != "" {
		req.Header.Set("Authorization", a.AuthToken)
	}
}
