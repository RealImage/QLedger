package middlewares

import (
	"net/http"

	ledgerContext "github.com/RealImage/QLedger/context"
)

type Handler func(http.ResponseWriter, *http.Request, *ledgerContext.AppContext)

func ContextMiddleware(handler Handler, context *ledgerContext.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, context)
	}
}
