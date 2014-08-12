package main

import (
	"errors"
	"net/http"
)

func postGame(r *http.Request) (*Game, *ApplicationError) {
	r.ParseForm()
	username := r.FormValue("username")
	if username == "" {
		msg := "Missing Parameter: username."
		err := errors.New("Missing Parameter")
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}
	game_name := r.FormValue("game_name")
	if game_name == "" {
		msg := "Missing Parameter: game_name."
		err := errors.New("Missing Parameter")
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}
	user, _ := GetUserByUsername(username)

	return NewGame(game_name, user.User_id)
}

func getGame(r *http.Request) ([]Game, *ApplicationError) {
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
