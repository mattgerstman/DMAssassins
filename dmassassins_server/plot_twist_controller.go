package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"errors"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	"net/http"
)

type PlotTwistPut struct {
	PlotTwistName  string `json:"plot_twist_name"`
	PlotTwistValue string `json:"plot_twist_value"`
	SendEmail      bool   `json:"send_email"`
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
		msg := "Invalid JSON"
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidJSON)
	}
	// Validate Name
	plotTwistName := plotTwistPut.PlotTwistName
	if plotTwistName == "" {
		msg := "Missing Parameter: plot_twist_name"
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}

	// Activate plot twist
	appErr = game.ActivatePlotTwist(plotTwistName)
	if appErr != nil {
		return nil, appErr
	}

	if !plotTwistPut.SendEmail {
		return game, nil
	}

	_, appErr = game.SendPlotTwistEmail(plotTwistName)
	if appErr != nil {
		LogWithSentry(appErr, map[string]string{"plot_twist_name": plotTwistName, "game_id": gameId.String()}, raven.WARNING)
	}

	return game, nil
}

// Handler for /game path
func GamePlotTwistHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "POST":
			fallthrough
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
