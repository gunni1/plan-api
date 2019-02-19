package api

import (
	"net/http"
)

type Route struct {
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

// Routes defines all HTTP routes, hanging off the main Server struct.
// Like that, all routes have access to the Server's dependencies.
func (s *Server) Routes() {
	var routes = Routes{
		Route{
			Method:      "GET",
			Pattern:     "/plan/{planId}",
			HandlerFunc: s.GetPlan(),
		},
		Route{
			Method:      "POST",
			Pattern:     "/plan",
			HandlerFunc: s.SavePlan(),
		},
		Route{
			Method:      "PUT",
			Pattern:     "/plan/{planId}",
			HandlerFunc: s.UpdatePlan(),
		},
	}

	for _, route := range routes {
		var handler http.HandlerFunc
		handler = route.HandlerFunc
		handler = Logger(handler)

		s.Router.
			Methods(route.Method).
			Path(route.Pattern).
			Handler(handler)

	}
}
