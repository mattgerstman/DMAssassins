package main

import (
	"database/sql"
	"encoding/json"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

var db *sql.DB

const (
	gameIdPath          = "/game/{game_id}/"
	gameLeaderboardPath = "/game/{game_id}/leaderboard/"
	gameUsersPath       = "/game/{game_id}/users/"
	gameUserPath        = "/game/{game_id}/user/{user_id}/"
	gameUserTargetPath  = "/game/{game_id}/user/{user_id}/target/"
	gameUserTeamPath    = "/game/{game_id}/user/{user_id}/team/{team_id}/"
	gameTeamPath        = "/game/{game_id}/team/"
	gameTeamIdPath      = "/game/{game_id}/team/{team_id}"
	gameRulesPath       = "/game/{game_id}/rules/"

	userGamePath = "/user/{user_id}/game/"
	sessionPath  = "/session/"
	homePath     = "/"

	HttpReponseCodeOk        = 200
	HttpResponseCodeCreated  = 201
	HttpReponseCodeNoContent = 204
)

// This function logs an error to the HTTP response and then returns an application error to be used as necessary
func HttpErrorLogger(w http.ResponseWriter, msg string, code int) {
	httpCode := code / 100
	http.Error(w, msg, httpCode)
}

// All HTTP requests should end up here, this function prints either an object or an error depending on the situation
// It also logs errors to sentry with a stack trace.
func WriteObjToPayload(w http.ResponseWriter, r *http.Request, obj interface{}, appErr *ApplicationError) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Content-Type", "application/json")
	if appErr != nil {
		HttpErrorLogger(w, appErr.Msg, appErr.Code)
		LogWithSentry(appErr, nil, raven.ERROR, raven.NewHttp(r))
		return
	}

	httpCode := HttpReponseCodeOk

	if obj == nil {
		httpCode = HttpReponseCodeNoContent
		w.Write(nil)
	}

	if (r.Method == "PUT") || (r.Method == "POST") {
		httpCode = HttpResponseCodeCreated
	}

	data, err := json.Marshal(obj)
	if err != nil {
		appErr := NewApplicationError("Internal Error", err, ErrCodeInternalServerWTF)
		LogWithSentry(appErr, nil, raven.ERROR, raven.NewHttp(r))
		HttpErrorLogger(w, appErr.Msg, appErr.Code)
		return
	}
	w.WriteHeader(httpCode)
	_, err = w.Write(data)
	if err != nil {
		appErr := NewApplicationError("Internal Error", err, ErrCodeInternalServerWTF)
		LogWithSentry(appErr, nil, raven.ERROR, raven.NewHttp(r))
		HttpErrorLogger(w, appErr.Msg, appErr.Code)
		return
	}
}

// Connects to the database, needs to be updated to read from an ini file
func connect() (db *sql.DB, err error) {
	db, err = sql.Open("postgres", Config.DatabaseURL)
	return db, err
}

// Handles CORS, eventually I'll strip it down to exactly the headers/origins I need
func corsHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Request-Headers", "X-Requested-With, accept, content-type")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, X-DMAssassins-Secret, X-DMAssassins-Game-Password, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		} else {
			h.ServeHTTP(w, r)
		}
	}
}

// Starts the server, opens the database, and registers handlers
func StartServer() {
	var err error
	db, err = connect()
	if err != nil {
		appErr := NewApplicationError("Could not connect to database", err, ErrCodeDatabase)
		LogWithSentry(appErr, nil, raven.ERROR)
		log.Fatal("Could not connect to database")
	}

	defer db.Close()

	r := mux.NewRouter().StrictSlash(true)

	// Just Game
	r.HandleFunc(gameIdPath, GameIdHandler()).Methods("POST", "GET", "DELETE")
	r.HandleFunc(gameLeaderboardPath, LeaderboardHandler()).Methods("GET")
	r.HandleFunc(gameRulesPath, GameRulesHandler()).Methods("GET", "POST")

	// Game then User
	r.HandleFunc(gameUserPath, GameUserHandler()).Methods("GET", "DELETE", "PUT")
	r.HandleFunc(gameUsersPath, GameUsersHandler()).Methods("GET", "DELETE", "PUT")
	r.HandleFunc(gameUserTargetPath, TargetHandler()).Methods("GET", "POST", "DELETE")
	r.HandleFunc(gameUserTeamPath, GameUserTeamHandler()).Methods("GET", "PUT", "POST", "DELETE")

	// Game then Team
	r.HandleFunc(gameTeamPath, GameTeamHandler()).Methods("GET", "POST")
	r.HandleFunc(gameTeamIdPath, GameTeamIdHandler()).Methods("GET", "POST", "DELETE")

	// User then Game
	r.HandleFunc(userGamePath, UserGameHandler()).Methods("GET", "PUT")

	// Just Session
	r.HandleFunc(sessionPath, SessionHandler()).Methods("POST")

	http.Handle("/", corsHandler(r))
	http.ListenAndServe(":8000", nil)
}
