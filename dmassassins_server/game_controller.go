package main

import (
	"errors"
	"net/http"
	"code.google.com/p/go-uuid/uuid"
)

func postGame(r *http.Request) (*Game, *ApplicationError) {
	r.ParseForm()
	userId := uuid.Parse(r.FormValue("user_id"))
	if userId == nil {
		msg := "Missing Parameter: user_id."
		err := errors.New("Missing Parameter")
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}
	gameName := r.FormValue("game_name")
	if gameName == "" {
		msg := "Missing Parameter: game_name."
		err := errors.New("Missing Parameter")
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}
	return NewGame(gameName, userId)
}

func getGame(r *http.Request) ([]*Game, *ApplicationError) {
	return GetGameList()
}

// Handler for /game path
func GameHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "GET":
			obj, err = getGame(r)

		case "POST":
			obj, err = postGame(r)
		default:
			obj = nil
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeInvalidMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
