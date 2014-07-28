package main

import (
	"database/sql"
	"encoding/json"
	//"errors"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	//"github.com/gorilla/schema"
)

var db *sql.DB

const (
	usersUsernamePath            = "/users/{username}/"
	usersUsernameTargetPath      = "/users/{username}/target/"
	usersUsernamePropertyPath    = "/users/{username}/property/"
	usersUsernamePropertyKeyPath = "/users/{username}/property/{key}/"
	sessionPath                  = "/session/"
	gamePath                     = "/game/"
	homePath                     = "/"
)

// This function logs an error to the HTTP response and then returns an application error to be used as necessary
func HttpErrorLogger(w http.ResponseWriter, msg string, code int) {
	httpCode := code / 100
	http.Error(w, msg, httpCode)
}

// All HTTP requests should end up here, this function prints either an object or an error depending on the situation
// It also logs errors to sentry with a stack trace.
func WriteObjToPayload(w http.ResponseWriter, r *http.Request, obj interface{}, err *ApplicationError) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println("Real Error\n") //debug line so I know errors I send vs ones from malformed paths
		fmt.Println(err)
		HttpErrorLogger(w, err.Msg, err.Code)
		LogWithSentry(err, nil, raven.ERROR, raven.NewHttp(r))
		return
	}

	var output map[string]interface{}
	output = make(map[string]interface{})
	output["response"] = obj
	encoder := json.NewEncoder(w)
	encoder.Encode(output)
}

// Handles requests to the direct path, currently does nothing
func HomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteObjToPayload(w, r, nil, nil)
	}
}

// Connects to the database, needs to be updated to read from an ini file
func connect() (*sql.DB, error) {
	var err error
	db, err = sql.Open("postgres", "postgres://localhost?dbname=dmassassins&sslmode=disable")
	fmt.Println(err)
	return db, err
}

// Starts the server, opens the database, and registers handlers
func StartServer() {
	connect()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc(homePath, HomeHandler()).Methods("GET")
	r.HandleFunc(usersUsernamePath, UserHandler()).Methods("GET")
	r.HandleFunc(usersUsernameTargetPath, TargetHandler()).Methods("GET", "POST")
	r.HandleFunc(usersUsernamePropertyKeyPath, UserPropertyHandler()).Methods("GET", "POST")

	r.HandleFunc(sessionPath, SessionHandler()).Methods("POST")
	// r.HandleFunc(gamePath, GameHandler()).Methods("POST")
	// Fuck you Taylor, this will be used again
	http.Handle("/", r)
	http.ListenAndServe(":8000", nil)
}
