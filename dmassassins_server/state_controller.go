package main

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

func postState(r *http.Request) (*Game, *ApplicationError) {
	vars := mux.Vars(r)
	game_name := vars["game_name"]

	game, appErr := GetGameByName(game_name)
	if appErr != nil {
		return nil, appErr
	}

	appErr = game.Start()
	if appErr != nil {
		return nil, appErr
	}

	return game, nil
}

func getState(r *http.Request) (*Game, *ApplicationError) {
	vars := mux.Vars(r)
	game_name := vars["game_name"]

	game, appErr := GetGameByName(game_name)
	if appErr != nil {
		return nil, appErr
	}
	return game, nil
}

func deleteState(r *http.Request) (*Game, *ApplicationError) {
	vars := mux.Vars(r)
	game_name := vars["game_name"]

	game, appErr := GetGameByName(game_name)
	if appErr != nil {
		return nil, appErr
	}

	appErr = game.End()
	if appErr != nil {
		return nil, appErr
	}
	return game, nil

}

// Handler for /game path
func StateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "GET":
			obj, err = getState(r)

		case "POST":
			obj, err = postState(r)
		default:
			obj = nil
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeInvalidMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
