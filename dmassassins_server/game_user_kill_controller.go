package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	"net/http"
)

// POST - Kills a user for their assassin
func postGameUserKill(r *http.Request) (appErr *ApplicationError) {
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

	assassin, appErr := GetAssassin(gameMapping.UserId, gameMapping.GameId)
	if appErr != nil {
		return appErr
	}

	secret := gameMapping.Secret

	_, oldTargetId, appErr := assassin.KillTarget(gameMapping.GameId, secret, true)
	if appErr != nil {
		return appErr
	}

	extra := GetExtraDataFromRequest(r)

	oldTarget, appErr := GetUserById(oldTargetId)
	if appErr != nil {
		LogWithSentry(appErr, map[string]string{"old_target_id": oldTargetId.String(), "game_id": gameId.String()}, raven.WARNING, extra)
		return nil
	}
	game, appErr := GetGameById(gameId)
	if appErr != nil {
		LogWithSentry(appErr, map[string]string{"old_target_id": oldTargetId.String(), "game_id": gameId.String()}, raven.WARNING, extra)
		return nil
	}
	_, appErr = oldTarget.SendDeadEmail(game.GameName)
	if appErr != nil {
		LogWithSentry(appErr, map[string]string{"old_target_id": oldTargetId.String(), "game_id": gameId.String()}, raven.WARNING, extra)
	}

	return nil
}

// Handler for /game path
func GameUserKillHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "POST":
			err = postGameUserKill(r)

		default:
			obj = nil
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeNotFoundMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
