package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	"net/http"
)

type AdminBanUserPost struct {
	SendEmail bool `json:"send_email"`
}

// DELETE - Bans a user from a game
func deleteGameUserBan(r *http.Request) (appErr *ApplicationError) {
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

	// Check if the user wants to send an email, if not just return
	sendEmail := r.Header.Get("X-DMAssassins-Send-Email")
	if sendEmail == "" {
		msg := "Missing Header: X-DMAssassins-Send-Email."
		err := errors.New(msg)
		return NewApplicationError(msg, err, ErrCodeMissingHeader)
	}

	if sendEmail == "false" {
		return nil
	}

	// Get game name and send the banhammer email
	sentryRequest := raven.NewHttp(r)

	user, appErr := GetUserById(userId)
	if appErr != nil {
		extra := map[string]interface{}{"user_id": userId.String()}
		LogWithSentry(appErr, nil, raven.WARNING, extra, sentryRequest)
		return nil

	}
	sentryUser := NewSentryUser(user)
	game, appErr := GetGameById(gameId)
	if appErr != nil {
		LogWithSentry(appErr, nil, raven.WARNING, nil, sentryRequest, sentryUser)
		return nil
	}
	_, appErr = user.SendBanhammerEmail(game.GameName)
	if appErr != nil {
		LogWithSentry(appErr, nil, raven.WARNING, nil, sentryRequest, sentryUser)
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
