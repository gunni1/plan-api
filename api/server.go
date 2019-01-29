package api

import (
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

type Server struct {
	Router  *mux.Router
	Session *mgo.Session
}
