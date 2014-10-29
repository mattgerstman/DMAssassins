package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	"net/http"
)

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

	extra := GetExtraDataFromRequest(r)

	assassin, appErr := GetUserById(assassinId)
	if appErr != nil {
		LogWithSentry(appErr, map[string]string{"user_id": userId.String()}, raven.WARNING, extra)
		return nil
	}

	game, appErr := GetGameById(gameId)
	if appErr != nil {
		LogWithSentry(appErr, map[string]string{"game_id": gameId.String()}, raven.WARNING, extra)
		return nil
	}
	_, appErr = assassin.SendNewTargetEmail(game.GameName)
	if appErr != nil {
		LogWithSentry(appErr, map[string]string{"user_id": userId.String(), "game_id": gameId.String()}, raven.WARNING, extra)
	}

	user, appErr := GetUserById(userId)
	if appErr != nil {
		LogWithSentry(appErr, map[string]string{"user_id": userId.String(), "game_id": gameId.String()}, raven.WARNING, extra)
		return nil
	}

	_, appErr = user.SendReviveEmail(game.GameName)
	if appErr != nil {
		LogWithSentry(appErr, map[string]string{"user_id": userId.String(), "game_id": gameId.String()}, raven.WARNING, extra)
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
