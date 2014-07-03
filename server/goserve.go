package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	//"fmt"
	//"github.com/gorilla/schema"
)

var db *sql.DB

const (
	usersPath = "/users/"
	loginPath = "/login/"
	gamePath  = "/game/"
	homePath  = "/"
)

func WriteObjToPayload(w http.ResponseWriter, obj interface{}) {
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

func StartServer() {
	connect()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc(homePath, HomeHandler()).Methods("GET")
	r.HandleFunc(usersPath, UserHandler()).Methods("GET", "POST", "DELETE")
	r.HandleFunc(loginPath, LoginHandler()).Methods("GET", "POST", "DELETE")
	r.HandleFunc(gamePath, GameHandler()).Methods("GET", "POST", "DELETE")
	http.Handle("/", r)
	http.ListenAndServe(":8000", nil)
}
