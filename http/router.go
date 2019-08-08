package http

import (
	"log"

	"github.com/gorilla/mux"
)

// NewRouter creates a new Router and initialises it will all of the global generateRoutes
func NewRouter(rs []route, log *log.Logger) (*mux.Router, error) {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range rs {
		handler := Logger(log, route.appHandler, route.name)
		router.
			Methods(route.method).
			Path(route.pattern).
			Name(route.name).
			Handler(handler)
		log.Printf("Route registered: %+v", route)
	}
	return router, nil
}
