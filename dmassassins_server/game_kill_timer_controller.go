package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	"net/http"
)

// get data for a kill timer
func getKillTimer(r *http.Request) (timerResponse map[string]interface{}, appErr *ApplicationError) {
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
	game, appErr := GetGameById(gameId)
	if appErr != nil {
		return nil, appErr
	}

	// Get timestamps response
	executeTs, createTs, appErr := game.GetTimesForGame()
	if appErr != nil {
		return nil, appErr
	}
	timerResponse = make(map[string]interface{})
	timerResponse["game_id"] = gameId.String()
	timerResponse["execute_ts"] = executeTs
	timerResponse["create_ts"] = createTs

	return timerResponse, nil
}

// stop a kill timer
func deleteKillTimer(r *http.Request) (appErr *ApplicationError) {
	_, appErr = RequiresAdmin(r)
	if appErr != nil {
		return appErr
	}

	// Get Game Id
	vars := mux.Vars(r)
	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id " + vars["game_id"]
		err := errors.New(msg)
		return NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	// Get game
	game, appErr := GetGameById(gameId)
	if appErr != nil {
		return appErr
	}

	appErr = game.DeleteKillTimer()
	if appErr != nil {
		return appErr
	}

	// Check if the user wants to send an email, if not just return
	sendEmail := r.Header.Get("X-DMAssassins-Send-Email")
	if sendEmail == "false" {
		return nil
	}

	// Inform users the game has ended
	_, appErr = game.SendTimerDisabledEmail()
	if appErr != nil {
		extra := make(map[string]interface{})
		extra[`game_id`] = gameId
		LogWithSentry(appErr, map[string]string{"game_id": gameId.String()}, raven.WARNING, extra)
	}

	return nil
}

// Handler for /game/<game_id>/kill_timer path
func GameKillTimerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "GET":
			obj, err = getKillTimer(r)
		case "DELETE":
			obj, err = nil, deleteKillTimer(r)
		default:
			obj = nil
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeNotFoundMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
