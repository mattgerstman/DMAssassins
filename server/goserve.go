package main

import (
	"database/sql"
	"encoding/json"
	//"errors"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	"github.com/getsentry/raven-go"
	//"github.com/gorilla/schema"
)

var db *sql.DB

const (
	usersGetPath = "/users/{email}"
	usersPath    = "/users/"
	loginPath    = "/login/"
	gamePath     = "/game/"
	homePath     = "/"
)

//This function logs an error to the HTTP response and then returns an application error to be used as necessary
func HttpErrorLogger(w http.ResponseWriter, msg string, code int) {
	httpCode := code / 100
	http.Error(w, msg, httpCode)
}

func WriteObjToPayload(w http.ResponseWriter, r *http.Request, obj interface{}, err *ApplicationError) {

	if err != nil {
		fmt.Println("Real Error\n") //debug line so I know errors I send vs ones from malformed paths
		fmt.Println(err)
		HttpErrorLogger(w, err.Msg, err.Code)
		LogWithSentry(err, nil, raven.ERROR, raven.NewHttp(r))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var output map[string]interface{}
	output = make(map[string]interface{})
	output["response"] = obj
	encoder := json.NewEncoder(w)
	encoder.Encode(output)
}

func HomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func connect() (*sql.DB, error) {
	var err error
	db, err = sql.Open("postgres", "postgres://localhost?dbname=dmassassins&sslmode=disable")
	fmt.Println(err)
	return db, err
}

//Is this right?
func StartServer() {
	connect()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc(homePath, HomeHandler()).Methods("GET")
	r.HandleFunc(usersGetPath, UserHandler()).Methods("GET")
	r.HandleFunc(usersPath, UserHandler()).Methods("POST", "DELETE")
	r.HandleFunc(loginPath, SessionHandler()).Methods("POST", "DELETE")
	r.HandleFunc(gamePath, GameHandler()).Methods("POST")
	http.Handle("/", r)
	http.ListenAndServe(":8000", nil)
}
