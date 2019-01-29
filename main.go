package plan_api

import (
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"os"
)

const (
	dbEnv = "PLAN_DB_URL"
)

func main() {
	mgoUrl := os.Getenv(dbEnv)
	if mgoUrl == "" {
		log.Fatalf("Environment variable %s is empty", dbEnv)
	}

	session, err := mgo.Dial(mgoUrl)
	if err != nil {
		log.Fatalf("Unable to establish DB connection to %s: %s", mgoUrl, err)
	}
	defer session.Close()
	log.Printf("DB connection established at %s", mgoUrl)

	server := api.Server{
		Router:  mux.NewRouter(),
		Session: session,
	}

}

func index(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(response, "Hello")
}
