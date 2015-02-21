package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
)

// GET function for /user/{user_id}/target returns a user's information
// Need to add permissions to this on a per user basis
func getTarget(r *http.Request) (user *User, appErr *ApplicationError) {
	_, appErr = RequiresUser(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)
	userId := uuid.Parse(vars["user_id"])
	if userId == nil {
		msg := "Invalid UUID: user_id " + vars["user_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}
	user, err := GetUserById(userId)
	if err != nil {
		return nil, err
	}
	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id " + vars["game_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	target, appErr := user.GetTarget(gameId)
	if appErr != nil {
		return nil, appErr
	}
	target.GetTeamByGameId(gameId)
	return target, nil
}

// DELETE - Kill a target, delete User may eventually be used by an admin
func deleteTarget(r *http.Request) (targetId uuid.UUID, appErr *ApplicationError) {
	_, appErr = RequiresUser(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)
	userId := uuid.Parse(vars["user_id"])
	if userId == nil {
		msg := "Invalid UUID: user_id " + vars["user_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	secret := r.Header.Get("X-DMAssassins-Secret")
	if secret == "" {
		msg := "Missing Header: X-DMAssassins-Secret."
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeMissingHeader)
	}

	user, err := GetUserById(userId)
	if err != nil {
		return nil, err
	}
	gameId := uuid.Parse(vars["game_id"])
	targetId, oldTargetId, appErr := user.KillTarget(gameId, secret, true)
	if appErr != nil {
		return nil, appErr
	}

	extra := map[string]interface{}{"old_target_id": oldTargetId.String(), "game_id": gameId.String()}
	sentryRequest := raven.NewHttp(r)
	sentryUser := NewSentryUser(user)

	oldTarget, appErr := GetUserById(oldTargetId)
	if appErr != nil {
		LogWithSentry(appErr, nil, raven.WARNING, extra, sentryRequest, sentryUser)
		return targetId, nil
	}
	game, appErr := GetGameById(gameId)
	if appErr != nil {
		LogWithSentry(appErr, nil, raven.WARNING, extra, sentryRequest, sentryUser)
		return targetId, nil
	}
	_, appErr = oldTarget.SendDeadEmail(game.GameName)
	if appErr != nil {
		LogWithSentry(appErr, nil, raven.WARNING, extra, sentryRequest, sentryUser)
	}

	return targetId, nil
}

// Handler for /user/{user_id}/target
func TargetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "GET":
			obj, err = getTarget(r)
		case "DELETE":
			obj, err = deleteTarget(r)
		default:
			obj = nil
			msg := "Not Found"
			err := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, err, ErrCodeNotFoundMethod)

		}
		WriteObjToPayload(w, r, obj, err)
	}
}
