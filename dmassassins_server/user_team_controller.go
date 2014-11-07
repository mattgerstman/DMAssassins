package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

// PUT - Add a user to a team
func putGameUserTeam(r *http.Request) (gameMapping *GameMapping, appErr *ApplicationError) {
	_, appErr = RequiresCaptain(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)
	userId := uuid.Parse(vars["user_id"])
	if userId == nil {
		msg := "Invalid UUID: user_id " + vars["user_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}

	// Get the user obj
	user, appErr := GetUserById(userId)
	if appErr != nil {
		return nil, appErr
	}

	teamId := uuid.Parse(vars["team_id"])
	if teamId == nil {
		msg := "Invalid UUID: team_id " + vars["team_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}
	return user.JoinTeam(teamId)
}

// GET - gets a team for a user by game_id
func getGameUserTeam(r *http.Request) (*Team, *ApplicationError) {
	_, appErr := RequiresCaptain(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)

	userId := uuid.Parse(vars["user_id"])
	if userId == nil {
		msg := "Invalid UUID: user_id " + vars["user_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}

	user, appErr := GetUserById(userId)
	if appErr != nil {
		return nil, appErr
	}

	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id " + vars["game_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}

	return user.GetTeamByGameId(gameId)
}

// DELETE - removes a user from a team
func deleteGameUserTeam(r *http.Request) (gameMapping *GameMapping, appErr *ApplicationError) {
	_, appErr = RequiresAdmin(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)

	userId := uuid.Parse(vars["user_id"])
	if userId == nil {
		msg := "Invalid UUID: user_id " + vars["user_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}

	// Get the user obj
	user, appErr := GetUserById(userId)
	if appErr != nil {
		return nil, appErr
	}

	teamId := uuid.Parse(vars["team_id"])
	if teamId == nil {
		msg := "Invalid UUID: team_id " + vars["team_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}
	return user.LeaveTeam(teamId)
}

// Handler for /team path
func GameUserTeamHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "GET":
			obj, err = getGameUserTeam(r)
		case "PUT":
			obj, err = putGameUserTeam(r)
		case "POST":
			obj, err = putGameUserTeam(r)
		case "DELETE":
			obj, err = deleteGameUserTeam(r)
		default:
			obj = nil
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeNotFoundMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
