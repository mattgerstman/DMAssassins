package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type PlotTwistPut struct {
	PlotTwistName  string `json:"plot_twist_name"`
	PlotTwistValue string `json:"plot_twist_value"`
}

func putPlotTwist(r *http.Request) (game *Game, appErr *ApplicationError) {
	//_, appErr = RequiresAdmin(r)
	if appErr != nil {
		return nil, appErr
	}

	// Get Game Id
	vars := mux.Vars(r)
	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id " + vars["game_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	// Get game
	game, appErr = GetGameById(gameId)
	if appErr != nil {
		return nil, appErr
	}

	// Decode json
	decoder := json.NewDecoder(r.Body)
	var plotTwistPut PlotTwistPut
	err := decoder.Decode(&plotTwistPut)
	if err != nil {
		return nil, NewApplicationError("Invalid JSON", err, ErrCodeInvalidJSON)
	}
	// Validate Name
	plotTwistName := plotTwistPut.PlotTwistName
	if plotTwistName == "" {
		return nil, NewApplicationError("Missing Parameter: plot_twist_name", err, ErrCodeMissingParameter)
	}
	plotTwistValue := plotTwistPut.PlotTwistValue
	if plotTwistValue == "" {
		return nil, NewApplicationError("Missing Parameter: plot_twist_value", err, ErrCodeMissingParameter)
	}

	// Activate plot twist
	appErr = game.ActivatePlotTwist(plotTwistName, plotTwistValue)
	if appErr != nil {
		return nil, appErr
	}

	return game, nil
}

// Handler for /game path
func GamePlotTwistHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "PUT":
			obj, err = putPlotTwist(r)
		default:
			obj = nil
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeNotFoundMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
