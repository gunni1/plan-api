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
	PLAN_DB_NAME              = "plan_db"
	PLAN_COLLECTION_NAME      = "plans"
	FAVORITES_COLLECTION_NAME = "favorites"
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
	CreatedBy string        `json:"createdBy" bson:"createdBy"`
	Practices []Practice    `json:"practices" bson:"practices"`
}

type UserFavorites struct {
	UserId        string   `bson:"_id,omitempty"`
	FavoritePlans []string `bson:"favoritePlans"`
}

type AddFavorite struct {
	PlanId string `json:"planId"`
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
		plans := s.Session.Copy().DB(PLAN_DB_NAME).C(PLAN_COLLECTION_NAME)

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
			SendErrorJSON(http.StatusBadRequest, decodeErr.Error(), w)
			return
		}
		planObjectId := bson.ObjectIdHex(planId)
		plans := s.Session.Copy().DB(PLAN_DB_NAME).C(PLAN_COLLECTION_NAME)

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

// Gets all Plans created by a given user. Sorted by title.
func (s *Server) GetUserPlans() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parameters := mux.Vars(r)
		userId := parameters["userId"]

		plans := s.Session.Copy().DB(PLAN_DB_NAME).C(PLAN_COLLECTION_NAME)
		userPlans := []Plan{}
		findErr := plans.Find(bson.M{"createdBy": userId}).Sort("title").All(&userPlans)
		if findErr != nil {
			SendErrorJSON(http.StatusInternalServerError, findErr.Error(), w)
		} else {
			NewResponse(http.StatusOK, userPlans).SendJSON(w)
		}
	}
}

// Gets the List of all favorite plans added by a given user.
// If the user has no favorites yet, returns a empty array
// Optional Parameter controls if only the PlanIds should given back
func (s *Server) GetUsersFavorites() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := mux.Vars(r)["userId"]
		onlyIds := r.URL.Query().Get("onlyIds")

		dbSession := s.Session.Copy()
		favorites := dbSession.DB(PLAN_DB_NAME).C(FAVORITES_COLLECTION_NAME)
		var userFavorites UserFavorites
		findErr := favorites.FindId(userId).One(&userFavorites)
		if findErr != nil {
			//If no Favorites created yet, return a empty array
			NewResponse(http.StatusOK, []UserFavorites{}).SendJSON(w)
		} else {
			if onlyIds == "true" {
				NewResponse(http.StatusOK, userFavorites.FavoritePlans).SendJSON(w)
			} else {
				//Get All plans of favorite Ids and Collect them as Result
				plans := dbSession.DB(PLAN_DB_NAME).C(PLAN_COLLECTION_NAME)
				planIds := make([]bson.ObjectId, len(userFavorites.FavoritePlans))
				for i := range planIds {
					planIds[i] = bson.ObjectIdHex(userFavorites.FavoritePlans[i])
				}
				favPlans := []Plan{}
				favPlanErr := plans.Find(bson.M{"_id": bson.M{"$in": planIds}}).All(&favPlans)
				if favPlanErr != nil {
					SendErrorJSON(http.StatusInternalServerError, favPlanErr.Error(), w)
				} else {
					NewResponse(http.StatusOK, favPlans).SendJSON(w)
				}
			}
		}
	}
}

// Add a planId to the List of users favorites
func (s *Server) AddFavorite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := mux.Vars(r)["userId"]
		var addFavDto AddFavorite
		decodeErr := json.NewDecoder(r.Body).Decode(&addFavDto)
		if decodeErr != nil {
			SendErrorJSON(http.StatusBadRequest, decodeErr.Error(), w)
			return
		}
		favorites := s.Session.Copy().DB(PLAN_DB_NAME).C(FAVORITES_COLLECTION_NAME)
		var userFavorites UserFavorites
		findErr := favorites.FindId(userId).One(&userFavorites)
		if findErr != nil {
			userFavorites = UserFavorites{
				UserId:        userId,
				FavoritePlans: []string{addFavDto.PlanId},
			}
			insertErr := favorites.Insert(&userFavorites)
			if insertErr != nil {
				SendErrorJSON(http.StatusInternalServerError, insertErr.Error(), w)
			} else {
				NewResponse(http.StatusOK, userFavorites.FavoritePlans).SendJSON(w)
			}
		} else {
			userFavorites.FavoritePlans = append(userFavorites.FavoritePlans, addFavDto.PlanId)
			favorites.UpdateId(userId, userFavorites)
			NewResponse(http.StatusOK, userFavorites.FavoritePlans).SendJSON(w)
		}
	}
}

// Remove a planId of the List of users favorites
func (s *Server) DelFavorite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := mux.Vars(r)["userId"]
		planId := mux.Vars(r)["planId"]

		favorites := s.Session.Copy().DB(PLAN_DB_NAME).C(FAVORITES_COLLECTION_NAME)
		var userFavorites UserFavorites
		findErr := favorites.FindId(userId).One(&userFavorites)
		if findErr != nil {
			SendErrorJSON(http.StatusBadRequest, ERR_NO_FAVORITES_FOR_USER, w)
		} else {
			index := getIndex(userFavorites.FavoritePlans, planId)
			if index == -1 {
				SendErrorJSON(http.StatusBadRequest, ERR_FAVORITE_NOT_FOUND, w)
			} else {
				newFavorites := append(userFavorites.FavoritePlans[:index], userFavorites.FavoritePlans[index+1:]...)
				userFavorites.FavoritePlans = newFavorites
				favorites.UpdateId(userId, userFavorites)
				NewResponse(http.StatusOK, userFavorites.FavoritePlans).SendJSON(w)
			}
		}
	}
}

func getIndex(s []string, e string) int {
	for i, v := range s {
		if v == e {
			return i
		}
	}
	return -1
}

func hasPlan(plans *mgo.Collection, planId bson.ObjectId) bool {
	count, err := plans.FindId(planId).Count()
	if err != nil {
		return false
	} else {
		return count > 0
	}
}
