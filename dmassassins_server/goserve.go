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
	gamePath               = "/game/"
	gameStatePath          = "/game/{game_id}/"
	gameUsernamePath       = "/{game_id}/users/{user_id}/"
	gameUsernameTargetPath = "/{game_id}/users/{user_id}/target/"
	sessionPath            = "/session/"
	homePath               = "/"
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
	w.Header().Set("Access-Control-Allow-Methods", "*")
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
	db, err = sql.Open("postgres", Config.DatabaseURL)
	fmt.Println(err)
	return db, err
}

// Catch all, This will eventually return a 404 but right now I'm using it to get request information
func CatchAllHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Println("Catch All")
		fmt.Println(r)
		WriteObjToPayload(w, r, r, nil)
	}
}

// Handles CORS, eventually I'll strip it down to exactly the headers/origins I need

func corsHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			//fmt.Println("CORS")
			fmt.Println(r)
			w.Header().Set("Access-Control-Request-Headers", "X-Requested-With, accept, content-type")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, X-DMAssassins-Secret")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
			//handle preflight in here
			//fmt.Println(w)
		} else {
			h.ServeHTTP(w, r)
		}
	}
}

// Starts the server, opens the database, and registers handlers
func StartServer() {
	_, err := connect()
	if err != nil {
		// add error here
	}

	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc(homePath, HomeHandler()).Methods("GET")
	r.HandleFunc(gameUsernamePath, UserHandler()).Methods("GET")
	r.HandleFunc(gameUsernameTargetPath, TargetHandler()).Methods("GET", "POST", "DELETE")

	r.HandleFunc(sessionPath, SessionHandler()).Methods("POST")

	r.HandleFunc(gamePath, GameHandler()).Methods("POST", "GET")
	r.HandleFunc(gameStatePath, StateHandler()).Methods("POST", "GET", "DELETE")

	r.HandleFunc("/{path:.*}", CatchAllHandler())
	http.Handle("/", corsHandler(r))
	http.ListenAndServe(":8000", nil)
}
