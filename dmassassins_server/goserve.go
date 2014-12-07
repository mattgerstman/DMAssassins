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
	"time"
)

var db *sql.DB

const (
	gameIdPath                = "/game/{game_id}/"
	gameLeaderboardPath       = "/game/{game_id}/leaderboard/"
	gameUsersPath             = "/game/{game_id}/users/"
	gameUsersEmailPath        = "/game/{game_id}/users/email/"
	gamePlotTwistPath         = "/game/{game_id}/plot_twist/"
	gameUserBanPath           = "/game/{game_id}/user/{user_id}/ban/"
	gameUserKillPath          = "/game/{game_id}/user/{user_id}/kill/"
	gameUserRevivePath        = "/game/{game_id}/user/{user_id}/revive/"
	gameUserPath              = "/game/{game_id}/user/{user_id}/"
	gameUserEmailPath         = "/game/{game_id}/user/{user_id}/email/"
	gameUserRolePath          = "/game/{game_id}/user/{user_id}/role/"
	gameUserTargetPath        = "/game/{game_id}/user/{user_id}/target/"
	gameUserTargetFriendsPath = "/game/{game_id}/user/{user_id}/target/friends/"
	gameUserTeamPath          = "/game/{game_id}/user/{user_id}/team/{team_id}/"
	gameTeamPath              = "/game/{game_id}/team/"
	gameTeamIdPath            = "/game/{game_id}/team/{team_id}/"
	gameRulesPath             = "/game/{game_id}/rules/"
	gameTargetsPath           = "/game/{game_id}/targets/"

	userGamePath    = "/user/{user_id}/game/"
	unsubscribePath = "/unsubscribe/{user_id}"
	sessionPath     = "/session/"
	homePath        = "/"

	HttpReponseCodeOk        = 200
	HttpResponseCodeCreated  = 201
	HttpReponseCodeNoContent = 204
)

// This function logs an error to the HTTP response and then returns an application error to be used as necessary
func HttpErrorLogger(w http.ResponseWriter, msg string, code int) {
	httpCode := code / 100
	http.Error(w, msg, httpCode)
}

// If we just want to return a string do it through this function
func WriteStringToPayload(w http.ResponseWriter, r *http.Request, msg string, appErr *ApplicationError) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Content-Type", "application/json")
	if appErr != nil {
		HttpErrorLogger(w, appErr.Msg, appErr.Code)
		extra := GetExtraDataFromRequest(r)
		LogWithSentry(appErr, nil, raven.ERROR, extra)
		return
	}
	httpCode := HttpReponseCodeOk
	w.WriteHeader(httpCode)
	byteMsg := []byte(msg)
	_, err := w.Write(byteMsg)
	if err != nil {
		appErr := NewApplicationError("Internal Error", err, ErrCodeInternalServerWTF)
		extra := GetExtraDataFromRequest(r)
		LogWithSentry(appErr, nil, raven.ERROR, extra)
		HttpErrorLogger(w, appErr.Msg, appErr.Code)
		return
	}

}

// All HTTP requests should end up here, this function prints either an object or an error depending on the situation
// It also logs errors to sentry with a stack trace.
func WriteObjToPayload(w http.ResponseWriter, r *http.Request, obj interface{}, appErr *ApplicationError) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Content-Type", "application/json")
	if appErr != nil {
		HttpErrorLogger(w, appErr.Msg, appErr.Code)
		extra := GetExtraDataFromRequest(r)
		LogWithSentry(appErr, nil, raven.ERROR, extra)
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
		extra := GetExtraDataFromRequest(r)
		LogWithSentry(appErr, nil, raven.ERROR, extra)
		return
	}
	w.WriteHeader(httpCode)
	_, err = w.Write(data)
	if err != nil {
		appErr := NewApplicationError("Internal Error", err, ErrCodeInternalServerWTF)
		extra := GetExtraDataFromRequest(r)
		LogWithSentry(appErr, nil, raven.ERROR, extra)
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
		// log.Println(r)
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Request-Headers", "X-Requested-With, accept, content-type")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, X-DMAssassins-Secret, X-DMAssassins-Game-Password, X-DMAssassins-Team-Id, Authorization")
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
		LogWithSentry(appErr, nil, raven.ERROR, nil)
		log.Fatal("Could not connect to database")
	}

	appErr := LoadAllTimers()
	if appErr != nil {
		fmt.Println(appErr)
		LogWithSentry(appErr, nil, raven.ERROR, nil)
	}

	defer db.Close()

	r := mux.NewRouter().StrictSlash(true)

	// Just Game
	r.HandleFunc(gameIdPath, GameIdHandler()).Methods("POST", "PUT", "GET", "DELETE")
	r.HandleFunc(gameLeaderboardPath, LeaderboardHandler()).Methods("GET")
	r.HandleFunc(gameRulesPath, GameRulesHandler()).Methods("GET", "POST")
	r.HandleFunc(gamePlotTwistPath, GamePlotTwistHandler()).Methods("PUT", "POST")
	r.HandleFunc(gameTargetsPath, GameTargetsHandler()).Methods("GET")

	// Game then User
	r.HandleFunc(gameUserPath, GameUserHandler()).Methods("GET", "DELETE", "PUT", "POST")
	r.HandleFunc(gameUsersPath, GameUsersHandler()).Methods("GET")
	r.HandleFunc(gameUsersEmailPath, GameUsersEmailHandler()).Methods("GET")
	r.HandleFunc(gameUserTargetPath, TargetHandler()).Methods("GET", "POST", "DELETE")
	r.HandleFunc(gameUserTargetFriendsPath, TargetFriendsHandler()).Methods("GET")
	r.HandleFunc(gameUserTeamPath, GameUserTeamHandler()).Methods("GET", "PUT", "POST", "DELETE")
	r.HandleFunc(gameUserRolePath, GameUserRoleHandler()).Methods("POST")

	// User actions
	r.HandleFunc(gameUserBanPath, GameUserBanHandler()).Methods("DELETE")
	r.HandleFunc(gameUserKillPath, GameUserKillHandler()).Methods("POST")
	r.HandleFunc(gameUserRevivePath, GameUserReviveHandler()).Methods("POST")

	// User Email Actions
	r.HandleFunc(gameUserEmailPath, GameUserEmailHandler()).Methods("POST")
	r.HandleFunc(unsubscribePath, UnsubscribeHandler()).Methods("GET")

	// Game then Team
	r.HandleFunc(gameTeamPath, GameTeamHandler()).Methods("GET", "POST")
	r.HandleFunc(gameTeamIdPath, GameTeamIdHandler()).Methods("GET", "POST", "DELETE", "PUT")

	// User then Game
	r.HandleFunc(userGamePath, UserGameHandler()).Methods("GET", "PUT")

	// Just Session
	r.HandleFunc(sessionPath, SessionHandler()).Methods("POST")

	http.Handle("/", corsHandler(http.TimeoutHandler(r, time.Second*20, `Timeout Occurred!`)))
	http.ListenAndServe(":8000", nil)
}
