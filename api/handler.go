package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"regexp"
)

const (
	PLAN_DB_NAME         = "plan_db"
	PLAN_COLLECTION_NAME = "plans"
)

var (
	objectIdRegEx, _ = regexp.Compile("[0-9a-fA-F]{24}")
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
		if !objectIdRegEx.MatchString(planId) {
			SendErrorJSON(http.StatusBadRequest, ERR_INVALID_PLAN_ID, w)
			return
		}
		var planDto Plan
		plans := s.Session.Copy().DB(PLAN_DB_NAME).C(PLAN_COLLECTION_NAME)
		err := plans.FindId(bson.ObjectIdHex(planId)).One(&planDto)
		if err != nil {
			SendErrorJSON(http.StatusBadRequest, ERR_NO_PLAN_FOUND, w)
			return
		}
		NewResponse(http.StatusOK, planDto).SendJSON(w)
		return
	}
}

// Saves a a given Json as a Plan
// Given ID in JSON is ignored. Always generating a new ObjectId
func (s *Server) SavePlan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestDto Plan
		decodeErr := json.NewDecoder(r.Body).Decode(&requestDto)
		if decodeErr != nil {
			SendErrorJSON(http.StatusBadRequest, decodeErr.Error(), w)
			return
		}
		dbSession := s.Session.Copy()
		plans := dbSession.DB(PLAN_DB_NAME).C(PLAN_COLLECTION_NAME)

		requestDto.Id = bson.NewObjectId()
		saveError := plans.Insert(&requestDto)
		if saveError != nil {
			SendErrorJSON(http.StatusBadRequest, ERR_NO_PLAN_FOUND, w)
		} else {
			NewResponse(http.StatusOK, requestDto).SendJSON(w)
		}
		return
	}
}

// Overrides a existing Plan by given Json (Full Update)
func (s *Server) UpdatePlan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parameters := mux.Vars(r)
		planId := parameters["planId"]
		if !objectIdRegEx.MatchString(planId) {
			SendErrorJSON(http.StatusBadRequest, ERR_INVALID_PLAN_ID, w)
			return
		}
		var requestDto Plan
		decodeErr := json.NewDecoder(r.Body).Decode(&requestDto)
		if decodeErr != nil {
			return
			SendErrorJSON(http.StatusBadRequest, decodeErr.Error(), w)
		}
		planObjectId := bson.ObjectIdHex(planId)
		dbSession := s.Session.Copy()
		plans := dbSession.DB(PLAN_DB_NAME).C(PLAN_COLLECTION_NAME)

		if hasPlan(plans, planObjectId) {
			requestDto.Id = planObjectId
			err := plans.UpdateId(planObjectId, &requestDto)
			if err != nil {
				SendErrorJSON(http.StatusBadRequest, err.Error(), w)
			} else {
				NewResponse(http.StatusOK, requestDto).SendJSON(w)
			}
		} else {
			SendErrorJSON(http.StatusBadRequest, ERR_NO_PLAN_FOUND, w)
		}
	}
}

func hasPlan(plans *mgo.Collection, planId bson.ObjectId) bool {
	count, err := plans.FindId(planId).Count()
	if err != nil {
		return false
	} else {
		return count > 0
	}
}
