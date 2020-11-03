package http

import (
	"github.com/glynternet/pkg/log"

	"net/http"
	"time"
)

// LoggingHandler returns a simple Handler that wraps a http.handler and logs any incoming request
func LoggingHandler(logger log.Logger, inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		duration := time.Since(start)
		_ = logger.Log(
			log.Message("Request handled"),
			log.KV{
				K: "method",
				V: r.Method,
			},
			log.KV{
				K: "requestURI",
				V: r.RequestURI,
			},
			log.KV{
				K: "name",
				V: name,
			},
			log.KV{
				K: "duration",
				V: duration,
			},
		)
	})
}
