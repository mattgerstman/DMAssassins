package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

// GET - Wrapper for GameMapping:GetGameMapping, usually used for user_role or alive/kill status
func getGameUsers(r *http.Request) (users map[string]*User, appErr *ApplicationError) {
	_, appErr = RequiresCaptain(r)
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

	users, appErr = game.GetAllUsersForGame()
	if appErr != nil {
		return nil, appErr
	}

	return users, nil
}

// Handler for /game path
func GameUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "GET":
			obj, err = getGameUsers(r)

		default:
			obj = nil
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeNotFoundMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
