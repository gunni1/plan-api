package plan_api

import (
	"gopkg.in/mgo.v2"
	"log"
)

const (
	trainingDbName         = "planDB"
	trainingCollectionName = "plans"
)

var (
	db *mgo.Database
)

func OpenDbConnection(dbUrl string) mgo.Session {
	log.Println("connect to mongodb using url: " + dbUrl)
	session, error := mgo.Dial(dbUrl)
	if error != nil {
		panic(error)
	}
	session.SetMode(mgo.Monotonic, true)
	db = session.DB(trainingDbName)
	return *session
}
