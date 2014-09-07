package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

// POST - Update rules for a game
func postGameRules(r *http.Request) (success string, appErr *ApplicationError) {
	appErr = RequiresAdmin(r)
	if appErr != nil {
		return "", appErr
	}

	vars := mux.Vars(r)
	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id" + gameId.String()
		err := errors.New(msg)
		return "", NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	game, appErr := GetGameById(gameId)
	if appErr != nil {
		return "", appErr
	}

	rules := r.FormValue("rules")
	if rules == "" {
		msg := "Missing Parameter: rules"
		err := errors.New(msg)
		return "", NewApplicationError(msg, err, ErrCodeMissingParameter)
	}

	appErr = game.SetRules(rules)
	if appErr != nil {
		return "", appErr
	}
	return "success", nil
}

// GET - Gets rules for a game
func getGameRules(r *http.Request) (rules string, appErr *ApplicationError) {
	appErr = RequiresLogin(r)
	if appErr != nil {
		return "", appErr
	}

	vars := mux.Vars(r)
	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id" + gameId.String()
		err := errors.New(msg)
		return "", NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	game, appErr := GetGameById(gameId)

	if appErr != nil {
		return "", appErr
	}

	return game.GetRules()
}

// Handler for /game path
func GameRulesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "GET":
			obj, err = getGameRules(r)
		case "POST":
			obj, err = postGameRules(r)
		default:
			obj = nil
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeNotFoundMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
