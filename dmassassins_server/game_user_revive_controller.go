package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"errors"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	"net/http"
)

type AdminReviveUserPost struct {
	SendEmail bool `json:"send_email"`
}

// POST - Revives a user and places them in between an assassin target apir
func postGameUserRevive(r *http.Request) (appErr *ApplicationError) {
	_, appErr = RequiresAdmin(r)
	if appErr != nil {
		return appErr
	}
	vars := mux.Vars(r)
	userId := uuid.Parse(vars["user_id"])
	if userId == nil {
		msg := "Invalid UUID: user_id " + vars["user_id"]
		err := errors.New(msg)
		return NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id " + vars["game_id"]
		err := errors.New(msg)
		return NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	gameMapping, appErr := GetGameMapping(userId, gameId)
	if appErr != nil {
		return appErr
	}

	assassinId, _, appErr := gameMapping.Revive()
	if appErr != nil {
		return appErr
	}

	// Check if the user wants to send an email, if not just return
	decoder := json.NewDecoder(r.Body)
	var adminReviveUserPost AdminReviveUserPost
	err := decoder.Decode(&adminReviveUserPost)
	if err != nil {
		return NewApplicationError("Invalid JSON", err, ErrCodeInvalidJSON)
	}

	extra := map[string]interface{}{"user_id": userId.String(), "game_id": gameId.String()}
	sentryRequest := raven.NewHttp(r)

	assassin, appErr := GetUserById(assassinId)
	if appErr != nil {
		LogWithSentry(appErr, nil, raven.WARNING, extra, sentryRequest)
		return nil
	}

	game, appErr := GetGameById(gameId)
	if appErr != nil {
		LogWithSentry(appErr, nil, raven.WARNING, extra, sentryRequest)
		return nil
	}
	_, appErr = assassin.SendNewTargetEmail(game.GameName)
	if appErr != nil {
		LogWithSentry(appErr, nil, raven.WARNING, extra, sentryRequest)
	}

	if !adminReviveUserPost.SendEmail {
		return nil
	}

	user, appErr := GetUserById(userId)
	if appErr != nil {
		LogWithSentry(appErr, nil, raven.WARNING, extra, sentryRequest)
		return nil
	}

	_, appErr = user.SendReviveEmail(game.GameName)
	if appErr != nil {
		LogWithSentry(appErr, nil, raven.WARNING, extra, sentryRequest)
	}

	return nil
}

// Handler for /game path
func GameUserReviveHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "POST":
			err = postGameUserRevive(r)

		default:
			obj = nil
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeNotFoundMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
