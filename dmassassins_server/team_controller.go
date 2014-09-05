package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

// POST - Creates a team for a game
func postGameTeam(r *http.Request) (team *Team, appErr *ApplicationError) {
	appErr = RequiresAdmin(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)
	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id" + gameId.String()
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidParameter)
	}
	game, appErr := GetGameById(gameId)
	if appErr != nil {
		return nil, appErr
	}

	teamName := r.FormValue("team_name")
	return game.NewTeam(teamName)
}

// GET - Gets the list of temas for a game
func getGameTeam(r *http.Request) (teams []*Team, appErr *ApplicationError) {
	appErr = RequiresAdmin(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)
	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id" + gameId.String()
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidParameter)
	}

	game, appErr := GetGameById(gameId)
	if appErr != nil {
		return nil, appErr
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
			err = NewApplicationError(msg, tempErr, ErrCodeInvalidMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
