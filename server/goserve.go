package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	"fmt"
	//"github.com/gorilla/schema"
)

var db *sql.DB

const (
	usersPath = "/users/"
	loginPath = "/login/"
	gamePath  = "/game/"
	homePath  = "/"
)
//This function logs an error to the HTTP response and then returns an application error to be used as necessary
func HttpErrorLogger(w http.ResponseWriter, msg string, code int) *ApplicationError {
	err := errors.New(msg)
	httpCode := code / 100
	http.Error(w, msg, httpCode)
	return &ApplicationError{msg, err, code}
}

func WriteObjToPayload(w http.ResponseWriter, obj interface{}, err *ApplicationError) {

	if err != nil {
		fmt.Println("Real Error\n") //debug line so I know errors I send vs ones from malformed paths
		HttpErrorLogger(w, err.Msg, err.Code)
		return
	}

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

func connect() {
	db, err = sql.Open("postgres", "postgres://localhost?dbname=dmassassins&sslmode=disable")
}

//Is this right?
func StartServer() {
	connect()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc(homePath, HomeHandler()).Methods("GET")
	r.HandleFunc(usersPath, UserHandler()).Methods("GET", "POST", "DELETE")
	r.HandleFunc(loginPath, SessionHandler()).Methods("POST", "DELETE")
	r.HandleFunc(gamePath, GameHandler()).Methods("GET", "POST", "DELETE")
	http.Handle("/", r)
	http.ListenAndServe(":8000", nil)
}
