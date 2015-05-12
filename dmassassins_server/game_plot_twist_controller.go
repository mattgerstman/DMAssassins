package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	"net/http"
)

type PlotTwistPost struct {
	PlotTwistName string `json:"plot_twist_name"`
	SendEmail     bool   `json:"send_email"`
}

// POST - creates a plot twist
func postPlotTwist(r *http.Request) (game *Game, appErr *ApplicationError) {
	_, appErr = RequiresAdmin(r)
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
	var plotTwistPost PlotTwistPost
	err := decoder.Decode(&plotTwistPost)
	if err != nil {
		msg := "Invalid JSON"
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidJSON)
	}
	// Validate Name
	plotTwistName := plotTwistPost.PlotTwistName
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

	fmt.Println(plotTwistPost)
	if !plotTwistPost.SendEmail {
		return game, nil
	}

	sentryRequest := raven.NewHttp(r)
	extra := map[string]interface{}{"plot_twist_name": plotTwistName, "game_id": gameId.String()}
	_, appErr = game.SendPlotTwistEmail(plotTwistName)
	if appErr != nil {
		LogWithSentry(appErr, nil, raven.WARNING, extra, sentryRequest)
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
			obj, err = postPlotTwist(r)
		default:
			obj = nil
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeNotFoundMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
