package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

const (
	PLAN_DB_NAME         = "plan_db"
	PLAN_COLLECTION_NAME = "plans"
)

type Practice struct {
	Name     string `json:"name" bson:"name"`
	Quantity string `json:"quantity" bson:"quantity"`
}

type Plan struct {
	Id        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Title     string        `json:"title" bson:"title"`
	Practices []Practice    `json:"practices" bson:"practices"`
}

// Retrieves a single Plan from the db
func (s *Server) GetPlan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parameters := mux.Vars(r)
		planId := parameters["planId"]

		var planDto Plan
		plans := s.Session.Copy().DB(PLAN_DB_NAME).C(PLAN_COLLECTION_NAME)
		err := plans.FindId(bson.ObjectIdHex(planId)).One(&planDto)
		if err != nil {
			log.Println(err)
			SendErrorJSON(http.StatusBadRequest, ERR_NO_PLAN_FOUND, w)
			return
		}
		NewResponse(http.StatusOK, planDto).SendJSON(w)
		return
	}
}

// Creates a Plan by given json
func (s *Server) CreatePlan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestDto Plan
		decodeErr := json.NewDecoder(r.Body).Decode(&requestDto)
		if decodeErr != nil {
			SendErrorJSON(http.StatusBadRequest, decodeErr.Error(), w)
			return
		}
		dbSession := s.Session.Copy()
		plans := dbSession.DB(PLAN_DB_NAME).C(PLAN_COLLECTION_NAME)

		dataDto := Plan{
			Id:        bson.NewObjectId(),
			Title:     requestDto.Title,
			Practices: requestDto.Practices,
		}
		err := plans.Insert(&dataDto)
		if err != nil {
			SendErrorJSON(http.StatusInternalServerError, err.Error(), w)
			return
		}

		NewResponse(http.StatusOK, dataDto).SendJSON(w)
		return
	}
}
