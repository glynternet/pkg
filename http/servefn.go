package http

import (
	"net/http"

	"github.com/glynternet/pkg/log"
)

// NewServeFn returns a function that can be used to start a server.
// NewServeFn will provide an HTTPS server if either the given certPath or
// keyPath are non-empty, otherwise NewServeFn will provide an HTTP server.
func NewServeFn(logger log.Logger, certPath, keyPath string) func(addr string, handler http.Handler) error {
	if len(certPath) == 0 && len(keyPath) == 0 {
		logger.Log(log.Message("Using HTTP"))
		return http.ListenAndServe
	}
	logger.Log(log.Message("Using HTTPS"))
	return func(addr string, handler http.Handler) error {
		return http.ListenAndServeTLS(addr, certPath, keyPath, handler)
	}
}
