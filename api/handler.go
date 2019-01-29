package api

import (
	"encoding/json"
	"net/http"
)

type Practice struct {
	Name string
	Reps string
}

type Plan struct {
	Title     string
	Practices []Practice
}

// Retrieves a single Plan from the db
func (s *Server) GetPlan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var planRequestDto Plan
		json.NewDecoder(r.Body).Decode(planRequestDto)
		return
	}
}

// Creates a Plan by given json
func (s *Server) CreatePlan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		return
	}
}
