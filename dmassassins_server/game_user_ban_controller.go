package main

import (
	"github.com/getsentry/raven-go"
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

// DELETE - Bans a user from a game
func deleteGameUserBan(r *http.Request) (appErr *ApplicationError) {
	_, appErr = RequiresAdmin(r)
	if appErr != nil {
		return appErr
	}
	vars := mux.Vars(r)
	userId := uuid.Parse(vars["user_id"])
	if userId == nil {
		msg := "Invalid UUID: user_id" + vars["user_id"]
		err := errors.New(msg)
		return NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id" + vars["game_id"]
		err := errors.New(msg)
		return NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	gameMapping, appErr := GetGameMapping(userId, gameId)
	if appErr != nil {
		return appErr
	}

	secret := gameMapping.Secret
	appErr = gameMapping.LeaveGame(secret)
	if appErr != nil {
		return appErr
	}

	user, appErr := GetUserById(userId)
	if appErr != nil {
		LogWithSentry(appErr, map[string]string{"user_id": userId.String()}, raven.WARNING)
		return nil

	}
	game, appErr := GetGameById(userId)
	if appErr != nil {
		LogWithSentry(appErr, map[string]string{"game_id": gameId.String()}, raven.WARNING)
		return nil
	}
	_, appErr = user.SendBanhammerEmail(game.GameName)
	if appErr != nil {
		LogWithSentry(appErr, map[string]string{"user_id": userId.String(), "game_id": gameId.String()}, raven.WARNING)
		return nil
	}

	return nil
}

// Handler for /game path
func GameUserBanHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "DELETE":
			err = deleteGameUserBan(r)

		default:
			obj = nil
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeNotFoundMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
