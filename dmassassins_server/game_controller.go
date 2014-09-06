package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"net/http"
)

// POST - Controller Wrapper for Game:NewGame
func postGame(r *http.Request) (game *Game, appErr *ApplicationError) {
	appErr = RequiresLogin(r)
	if appErr != nil {
		return nil, appErr
	}

	r.ParseForm()
	userId := uuid.Parse(r.FormValue("user_id"))
	if userId == nil {
		msg := "Invalid Parameter: user_id " + userId.String()
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	gameName := r.FormValue("game_name")
	if gameName == "" {
		msg := "Missing Parameter: game_name."
		err := errors.New("Missing Parameter")
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}

	gamePassword := r.FormValue("game_password")
	return NewGame(gameName, userId, gamePassword)
}

// GET - Controller wrapper for Game::GetGamesList
func getGame(r *http.Request) (games []*Game, appErr *ApplicationError) {
	appErr = RequiresLogin(r)
	if appErr != nil {
		return nil, appErr
	}
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
			err = NewApplicationError(msg, tempErr, ErrCodeNotFoundMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
