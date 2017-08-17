package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// TokenAuthMiddleware is a middleware that provides authentication functionality
func TokenAuthMiddleware(router *httprouter.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check whether token authentication enabled
		envToken := strings.TrimSpace(os.Getenv("LEDGER_AUTH_TOKEN"))
		if envToken != "" {
			// Get the token in the header
			requestToken := strings.TrimSpace(r.Header.Get("LEDGER-AUTH-TOKEN"))
			// Validate token
			if requestToken != envToken {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		router.ServeHTTP(w, r)
	}
}
