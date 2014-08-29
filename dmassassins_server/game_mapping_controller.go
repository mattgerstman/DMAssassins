package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

func postGameMapping(r *http.Request) (*GameMapping, *ApplicationError) {
	appErr := RequiresLogin(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)
	userId := uuid.Parse(vars["user_id"])
	gameId := uuid.Parse(vars["game_id"])
	gamePassword := r.FormValue("game_password")

	user, appErr := GetUserById(userId)
	if appErr != nil {
		return nil, appErr
	}

	return user.JoinGame(gameId, gamePassword)
}

func getGameMapping(r *http.Request) (*GameMapping, *ApplicationError) {
	appErr := RequiresUser(r)
	if appErr != nil {
		return nil, appErr
	}
	vars := mux.Vars(r)
	userId := uuid.Parse(vars["user_id"])
	gameId := uuid.Parse(vars["game_id"])
	return GetGameMapping(userId, gameId)
}

// Handler for /game path
func GameMappingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "GET":
			obj, err = getGameMapping(r)

		case "POST":
			obj, err = postGameMapping(r)
		default:
			obj = nil
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeInvalidMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
