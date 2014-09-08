package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

var db *sql.DB

const (
	gamePath            = "/game/"
	gameIdPath          = "/game/{game_id}/"
	gameLeaderboardPath = "/game/{game_id}/leaderboard/"
	gameUserPath        = "/game/{game_id}/users/{user_id}/"
	gameUserTargetPath  = "/game/{game_id}/users/{user_id}/target/"
	gameUserTeamPath    = "/game/{user_id}/users/{user_id}/team/"
	gameTeamPath        = "/game/{game_id}/team/"
	gameRulesPath       = "/game/{game_id}/rules/"
	teamIdPath          = "/team/{team_id}"
	userPath            = "/users/{user_id}/"
	userGamePath        = "/users/{user_id}/game/"
	userGameNewPath     = "/users/{user_id}/game/new/"
	sessionPath         = "/session/"
	homePath            = "/"
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

	var output map[string]interface{}
	output = make(map[string]interface{})
	output["response"] = obj

	data, err := json.Marshal(output)
	if err != nil {
		appErr := NewApplicationError("Internal Error", err, ErrCodeInternalServerWTF)
		LogWithSentry(appErr, nil, raven.ERROR, raven.NewHttp(r))
		HttpErrorLogger(w, appErr.Msg, appErr.Code)
		return
	}
	w.Write(data)
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
			fmt.Println(r)
			w.Header().Set("Access-Control-Request-Headers", "X-Requested-With, accept, content-type")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, X-DMAssassins-Secret, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
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

	r := mux.NewRouter()

	// Just Game
	r.HandleFunc(gamePath, GameHandler()).Methods("POST", "GET")
	r.HandleFunc(gameIdPath, GameIdHandler()).Methods("POST", "GET", "DELETE")
	r.HandleFunc(gameLeaderboardPath, LeaderboardHandler()).Methods("GET")
	r.HandleFunc(gameRulesPath, GameRulesHandler()).Methods("GET", "POST")

	// Game then User
	r.HandleFunc(gameUserPath, GameUserHandler()).Methods("GET", "POST")
	r.HandleFunc(gameUserTargetPath, TargetHandler()).Methods("GET", "POST", "DELETE")
	r.HandleFunc(gameUserTeamPath, GameUserTeamHandler()).Methods("GET", "POST")

	// Game then Team
	r.HandleFunc(gameTeamPath, GameTeamHandler()).Methods("GET", "POST")

	// Just User
	r.HandleFunc(userPath, UserHandler()).Methods("GET")

	// User then Game
	r.HandleFunc(userGamePath, UserGameHandler()).Methods("GET")
	r.HandleFunc(userGameNewPath, UserGameNewHandler()).Methods("GET")

	// Just Team
	r.HandleFunc(teamIdPath, TeamIdHandler()).Methods("GET", "POST", "DELETE")

	// Just Session
	r.HandleFunc(sessionPath, SessionHandler()).Methods("POST")

	http.Handle("/", corsHandler(r))
	http.ListenAndServe(":8000", nil)
}
