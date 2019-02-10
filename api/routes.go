package api

import (
	"github.com/rs/cors"
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
	}

	for _, route := range routes {
		var handler http.HandlerFunc
		handler = route.HandlerFunc
		handler = Logger(handler)

		//CORS Header
		c := cors.New(cors.Options{
			AllowedMethods:     []string{"GET", "POST", "OPTIONS"},
			AllowedOrigins:     []string{"http://localhost:4200"},
			AllowCredentials:   true,
			AllowedHeaders:     []string{"Content-Type", "Bearer", "Bearer ", "content-type", "Origin", "Accept"},
			OptionsPassthrough: true,
		})

		//Options explizit immer erlauben für Preflight-Check
		s.Router.
			Methods(route.Method, "OPTIONS").
			Path(route.Pattern).
			Handler(c.Handler(handler))

	}
}
