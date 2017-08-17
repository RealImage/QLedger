package middlewares

import (
	"net/http"

	ledgerContext "github.com/RealImage/QLedger/context"
)

// Handler is a custom HTTP handler that has an additional application context
type Handler func(http.ResponseWriter, *http.Request, *ledgerContext.AppContext)

// ContextMiddleware is a middleware that provides application context to the `Handler`
func ContextMiddleware(handler Handler, context *ledgerContext.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, context)
	}
}
