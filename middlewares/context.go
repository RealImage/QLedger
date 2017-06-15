package middlewares

import (
	"net/http"

	ledgerContext "github.com/RealImage/QLedger/context"
	"github.com/julienschmidt/httprouter"
)

type Handler func(http.ResponseWriter, *http.Request, *ledgerContext.AppContext)
type HandlerWithParams func(http.ResponseWriter, *http.Request, httprouter.Params, *ledgerContext.AppContext)

func ContextMiddleware(handler Handler, context *ledgerContext.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, context)
	}
}

func ContextParamsMiddleware(handler HandlerWithParams, context *ledgerContext.AppContext) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		handler(w, r, params, context)
	}
}
