package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type RulesPost struct {
	Rules string `json:"rules"`
}

// put - Update rules for a game
func putGameRules(r *http.Request) (appErr *ApplicationError) {
	_, appErr = RequiresAdmin(r)
	if appErr != nil {
		return appErr
	}

	vars := mux.Vars(r)
	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id " + vars["game_id"]
		err := errors.New(msg)
		return NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	game, appErr := GetGameById(gameId)
	if appErr != nil {
		return appErr
	}

	decoder := json.NewDecoder(r.Body)
	var rulesPost RulesPost
	err := decoder.Decode(&rulesPost)
	if err != nil {
		return NewApplicationError("Invalid JSON", err, ErrCodeInvalidJSON)
	}

	rules := rulesPost.Rules
	if rules == "" {
		msg := "Missing Parameter: rules"
		err := errors.New(msg)
		return NewApplicationError(msg, err, ErrCodeMissingParameter)
	}

	appErr = game.SetRules(rules)
	if appErr != nil {
		return appErr
	}
	return nil
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
		msg := "Invalid UUID: game_id " + vars["game_id"]
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
		case "PUT":
			obj, err = nil, putGameRules(r)
		default:
			obj = nil
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeNotFoundMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
