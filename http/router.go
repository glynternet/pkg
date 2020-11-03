package http

import (
	"github.com/glynternet/pkg/log"

	"github.com/gorilla/mux"
)

// NewRouter creates a new Router and initialises it will all of the global generateRoutes
func NewRouter(rs []Route, logger log.Logger) (*mux.Router, error) {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range rs {
		handler := LoggingHandler(logger, route.AppHandler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
		logger.Log(log.Message("Route registered"), log.KV{
			K: "route",
			V: route,
		})
	}
	return router, nil
}
