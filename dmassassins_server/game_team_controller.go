package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type TeamPost struct {
	TeamName string `json:"team_name"`
}

// POST - Creates a team for a game
func postGameTeam(r *http.Request) (team *Team, appErr *ApplicationError) {
	_, appErr = RequiresAdmin(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)
	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id " + vars["game_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}
	game, appErr := GetGameById(gameId)
	if appErr != nil {
		return nil, appErr
	}

	decoder := json.NewDecoder(r.Body)
	var teamInfo TeamPost
	err := decoder.Decode(&teamInfo)
	if err != nil {
		return nil, NewApplicationError("Invalid JSON", err, ErrCodeInvalidJSON)
	}

	teamName := teamInfo.TeamName
	if teamName == "" {
		msg := "Missing Parameter: team_name"
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}

	return game.NewTeam(teamName)
}

// GET - Gets the list of teams for a game
func getGameTeam(r *http.Request) (teams []*Team, appErr *ApplicationError) {
	appErr = RequiresLogin(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)
	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id " + vars["game_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	game, appErr := GetGameById(gameId)
	if appErr != nil {
		return nil, appErr
	}

	teamsEnabled, appErr := game.GetGameProperty(`teams_enabled`)
	if appErr != nil {
		return nil, appErr
	}

	if teamsEnabled != `true` {
		return nil, nil
	}

	return game.GetTeams()
}

// Handler for /team path
func GameTeamHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "GET":
			obj, err = getGameTeam(r)
		case "POST":
			obj, err = postGameTeam(r)
		default:
			obj = nil
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeNotFoundMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
