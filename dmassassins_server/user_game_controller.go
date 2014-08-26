package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

func postUserGame(r *http.Request) (*Game, *ApplicationError) {
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
	gameName := r.FormValue("game_name")
	if gameName == "" {
		msg := "Missing Parameter: game_name."
		err := errors.New("Missing Parameter")
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}
	return nil, nil
	//return NewGame(gameName, userId)
}

func getUserGame(r *http.Request) ([]*Game, *ApplicationError) {
	appErr := RequiresUser(r)
	if appErr != nil {
		return nil, appErr
	}
	vars := mux.Vars(r)
	userId := uuid.Parse(vars["user_id"])

	user, appErr := GetUserById(userId)
	if appErr != nil {
		return nil, appErr
	}

	return user.GetGamesForUser()
}

// Handler for /game path
func UserGameHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "GET":
			obj, err = getUserGame(r)

		case "POST":
			obj, err = postUserGame(r)
		default:
			obj = nil
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeInvalidMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
