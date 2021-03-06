package api

import (
	"net/http"
)

type Route struct {
	Method      string
	Pattern     string
	Queries     string
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
		Route{
			Method:      "DELETE",
			Pattern:     "/plan/{planId}",
			HandlerFunc: s.DeletePlan(),
		},
		Route{
			Method:      "GET",
			Pattern:     "/userplans/{userId}",
			HandlerFunc: s.GetUserPlans(),
		},
		Route{
			Method:      "GET",
			Pattern:     "/userfav/{userId}",
			HandlerFunc: s.GetUsersFavorites(),
		},
		Route{
			Method:      "POST",
			Pattern:     "/userfav/{userId}/plan",
			HandlerFunc: s.AddFavorite(),
		},
		Route{
			Method:      "DELETE",
			Pattern:     "/userfav/{userId}/plan/{planId}",
			HandlerFunc: s.DelFavorite(),
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
