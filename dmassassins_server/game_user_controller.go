package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"errors"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	"net/http"
)

type UserPropertyPost struct {
	Properties map[string]string `json:properties`
}

// PUT - Wrapper for UserProperties::SetUserProperty
func putGameUser(r *http.Request) (user *User, appErr *ApplicationError) {
	_, appErr = RequiresCaptain(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)
	userId := uuid.Parse(vars["user_id"])
	if userId == nil {
		msg := "Invalid UUID: user_id" + vars["user_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}
	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id" + vars["game_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	decoder := json.NewDecoder(r.Body)
	var newProperties UserPropertyPost
	err := decoder.Decode(&newProperties)
	if err != nil {
		return nil, NewApplicationError("Invalid JSON", err, ErrCodeInvalidJSON)
	}

	user, appErr = GetUserById(userId)
	if appErr != nil {
		return nil, appErr
	}

	allowedProperties := make(map[string]string)
	if _, ok := newProperties.Properties[`photo`]; !ok {
		return user, nil
	}

	newPhoto := newProperties.Properties[`photo`]

	oldPhoto, appErr := user.GetUserProperty(`photo`)
	if newPhoto == oldPhoto {
		return nil, nil
	}

	allowedProperties[`photo`] = newPhoto
	allowedProperties[`photo_thumb`] = newPhoto
	appErr = user.SetUserProperties(allowedProperties)
	if appErr != nil {
		return nil, appErr
	}

	return nil, nil
}

// POST - Wrapper for GameMapping:JoinGame
func postGameUser(r *http.Request) (gameMapping *GameMapping, appErr *ApplicationError) {
	appErr = RequiresLogin(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)
	userId := uuid.Parse(vars["user_id"])
	if userId == nil {
		msg := "Invalid UUID: user_id" + vars["user_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}
	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id" + vars["game_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	gamePassword := r.Header.Get("X-DMAssassins-Game-Password")

	user, appErr := GetUserById(userId)
	if appErr != nil {
		return nil, appErr
	}

	gameMapping, appErr = user.JoinGame(gameId, gamePassword)
	if appErr != nil {
		return nil, appErr
	}

	teamIdHeader := r.Header.Get("X-DMAssassins-Team-Id")
	teamId := uuid.Parse(teamIdHeader)
	if teamId == nil {
		return gameMapping, nil
	}

	extra := GetExtraDataFromRequest(r)

	user, appErr = GetUserById(userId)
	if appErr != nil {
		LogWithSentry(appErr, map[string]string{"user_id": userId.String(), "game_id": gameId.String(), "team_id": teamId.String()}, raven.WARNING, extra)
		return gameMapping, nil
	}
	gameMapping, appErr = user.JoinTeam(teamId)
	if appErr != nil {
		LogWithSentry(appErr, map[string]string{"user_id": userId.String(), "game_id": gameId.String(), "team_id": teamId.String()}, raven.WARNING, extra)
	}

	return gameMapping, nil
}

// GET - Wrapper for GameMapping:GetGameMapping, usually used for user_role or alive/kill status
func getGameUser(r *http.Request) (user *User, appErr *ApplicationError) {
	_, appErr = RequiresUser(r)
	if appErr != nil {
		return nil, appErr
	}
	vars := mux.Vars(r)
	userId := uuid.Parse(vars["user_id"])
	if userId == nil {
		msg := "Invalid UUID: user_id" + vars["user_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id" + vars["game_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	user, appErr = GetUserForGameById(userId, gameId)
	if appErr != nil {
		return nil, appErr
	}

	return user, nil
}

// DELETE - Lets a user quit the game
func deleteGameUser(r *http.Request) (appErr *ApplicationError) {
	_, appErr = RequiresUser(r)
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

	r.ParseForm()
	secret := r.Header.Get("X-DMAssassins-Secret")
	if secret == "" {
		msg := "Missing Header: X-DMAssassins-Secret."
		err := errors.New(msg)
		return NewApplicationError(msg, err, ErrCodeMissingHeader)
	}

	return gameMapping.LeaveGame(secret)
}

// Handler for /game path
func GameUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "GET":
			obj, err = getGameUser(r)
		case "PUT":
			obj, err = putGameUser(r)
		case "POST":
			obj, err = postGameUser(r)
		case "DELETE":
			err = deleteGameUser(r)

		default:
			obj = nil
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeNotFoundMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
