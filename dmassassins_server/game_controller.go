package main

import (
	"encoding/json"
	"errors"
	"github.com/getsentry/raven-go"
	"net/http"
)

type NewGamePost struct {
	GameName     string `json:"game_name"`
	GamePassword string `json:"game_password"`
}

// POST - Controller Wrapper for Game:NewGame
func postGame(r *http.Request) (game *Game, appErr *ApplicationError) {
	appErr = RequiresLogin(r)
	if appErr != nil {
		return nil, appErr
	}

	user := GetUserForRequest(r)
	if user == nil {
		msg := "Internal Error"
		err := errors.New("Missing user for request")
		return nil, NewApplicationError(msg, err, ErrCodeNoUserForContext)
	}

	decoder := json.NewDecoder(r.Body)
	var newGame NewGamePost
	err := decoder.Decode(&newGame)
	if err != nil {
		return nil, NewApplicationError("Invalid JSON", err, ErrCodeInvalidJSON)
	}

	gameName := newGame.GameName
	if gameName == "" {
		msg := "Missing Parameter: game_name."
		err := errors.New("Missing Parameter")
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}

	gamePassword := newGame.GamePassword
	game, appErr = NewGame(gameName, user.UserId, gamePassword)
	if appErr != nil {
		return nil, appErr
	}

	sentryRequest := raven.NewHttp(r)
	sentryUser := NewSentryUser(user)
	_, appErr = user.SendAdminWelcomeEmail()
	if appErr != nil {
		LogWithSentry(appErr, nil, raven.WARNING, nil, sentryRequest, sentryUser)
	}
	return game, nil
}

// GET - gets a list of games for a user
func getGame(r *http.Request) (response map[string][]*Game, appErr *ApplicationError) {
	appErr = RequiresLogin(r)
	if appErr != nil {
		return nil, appErr
	}

	user := GetUserForRequest(r)
	if user == nil {
		msg := "Internal Error"
		err := errors.New("Missing user for request")
		return nil, NewApplicationError(msg, err, ErrCodeNoUserForContext)
	}

	response = make(map[string][]*Game)
	response["member"], appErr = user.GetGamesForUser()
	if appErr != nil {
		return nil, appErr
	}
	response["available"], appErr = user.GetNewGamesForUser()
	if appErr != nil {
		return nil, appErr
	}
	return response, nil
}

// Handler for /game path
func GameHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "POST":
			obj, err = postGame(r)
		case "GET":
			obj, err = getGame(r)
		default:
			obj = nil
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeNotFoundMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
