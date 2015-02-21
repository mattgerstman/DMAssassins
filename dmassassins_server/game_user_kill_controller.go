package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"errors"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	"net/http"
)

type AdminKillUserPost struct {
	SendEmail bool `json:"send_email"`
}

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
	if appErr != nil && appErr.Code != ErrCodeNotFoundUserId {
		return appErr
	}

	// if we don't have an assassin just kill the damn user
	if appErr != nil {
		return gameMapping.MarkDead()
	}

	secret := gameMapping.Secret
	_, oldTargetId, appErr := assassin.KillTarget(gameMapping.GameId, secret, true)
	if appErr != nil {
		return appErr
	}

	// Check if the user wants to send an email, if not just return
	decoder := json.NewDecoder(r.Body)
	var adminKillUserPost AdminKillUserPost
	err := decoder.Decode(&adminKillUserPost)
	if err != nil {
		return NewApplicationError("Invalid JSON", err, ErrCodeInvalidJSON)
	}

	if !adminKillUserPost.SendEmail {
		return nil
	}

	// Get game name and send the banhammer email

	sentryRequest := raven.NewHttp(r)
	extra := map[string]interface{}{"assassin_id": assassin.UserId.String(), "old_target_id": oldTargetId.String(), "game_id": gameId.String()}

	oldTarget, appErr := GetUserById(oldTargetId)
	if appErr != nil {
		LogWithSentry(appErr, nil, raven.WARNING, extra, sentryRequest)
		return nil
	}
	game, appErr := GetGameById(gameId)
	if appErr != nil {
		LogWithSentry(appErr, nil, raven.WARNING, extra, sentryRequest)
		return nil
	}
	_, appErr = oldTarget.SendDeadEmail(game.GameName)
	if appErr != nil {
		LogWithSentry(appErr, nil, raven.WARNING, extra, sentryRequest)
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
