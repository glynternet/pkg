package http

import (
	"net/http"

	"github.com/glynternet/pkg/log"
)

// WithAuthoriser wraps a handler with the given RequestAuthoriser.
// When handling requests, if the request is not authorised, the next
// http.Handler is not called.
func WithAuthoriser(logger log.Logger, ra RequestAuthoriser, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := ra.Authorise(r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			if logger != nil {
				_ = logger.Log(
					log.Message("Unauthorised request"),
					log.ErrorMessage(err))
			}
			return
		}
		next.ServeHTTP(w, r)
	}
}

// RequestAuthoriser returns an error of the request is not authorised.
type RequestAuthoriser interface {
	Authorise(*http.Request) error
}
