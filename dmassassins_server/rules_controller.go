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
	return "", nil
}

// GET - Gets rules for a game
func getGameRules(r *http.Request) (rulesWrapper map[string]string, appErr *ApplicationError) {
	appErr = RequiresLogin(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)
	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id" + gameId.String()
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	game, appErr := GetGameById(gameId)
	if appErr != nil {
		return nil, appErr
	}

	rules, appErr := game.GetRules()
	if appErr != nil {
		return nil, appErr
	}

	rulesWrapper = make(map[string]string)
	rulesWrapper["rules"] = rules

	return rulesWrapper, nil
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
