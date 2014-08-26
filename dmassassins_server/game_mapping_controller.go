package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"net/http"
)

func postGameMapping(r *http.Request) (*GameMapping, *ApplicationError) {
	appErr := RequiresLogin(r)
	if appErr != nil {
		return nil, appErr
	}

	r.ParseForm()
	userId := uuid.Parse(r.FormValue("user_id"))
	if userId == nil {
		msg := "Missing Parameter: user_id."
		err := errors.New("Missing Parameter")
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}
	gameId := uuid.Parse(r.FormValue("game_id"))
	if gameId == nil {
		msg := "Missing Parameter: game_id."
		err := errors.New("Missing Parameter")
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}

	user, appErr := GetUserById(userId)
	if appErr != nil {
		return nil, appErr
	}

	return user.JoinGame(gameId)
}

func getGameMapping(r *http.Request) ([]*Game, *ApplicationError) {
	appErr := RequiresLogin(r)
	if appErr != nil {
		//return nil, appErr
	}
	return GetGameList()
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
