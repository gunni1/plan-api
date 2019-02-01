package api

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

const (
	PLAN_DB_NAME         = "plan_db"
	PLAN_COLLECTION_NAME = "plans"
)

type Practice struct {
	Name     string `json:"name" bson: "name"`
	Quantity string `json:"quantity" bson: "quantity"`
}

type Plan struct {
	Id        bson.ObjectId `json: "id" bson: "_id" `
	Title     string        `json:"title" bson: "title"`
	Practices []Practice    `json:"practices" bson: "practices"`
}

// Retrieves a single Plan from the db
func (s *Server) GetPlan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		NewResponse(http.StatusNotImplemented, "try later").SendJSON(w)

		return
	}
}

// Creates a Plan by given json
func (s *Server) CreatePlan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var planDto Plan
		decodeErr := json.NewDecoder(r.Body).Decode(planDto)
		if decodeErr != nil {
			http.Error(w, decodeErr.Error(), http.StatusBadRequest)
			return
		}

		dbSession := s.Session.Copy()
		plans := dbSession.DB(PLAN_DB_NAME).C(PLAN_COLLECTION_NAME)

		err := plans.Insert(planDto)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}
}
