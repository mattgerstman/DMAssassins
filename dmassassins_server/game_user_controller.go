package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"errors"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	"net/http"
)

type UserPost struct {
	Email      string            `json:"email"`
	Properties map[string]string `json:properties`
}

// Makes sure we only allow the user to set properties that are explicitly allowed
func filterProperties(properties map[string]string) (filteredProperties map[string]string) {
	filteredProperties = make(map[string]string)

	_, havePhoto := properties[`photo`]
	_, havePhotoThumb := properties[`photo_thumb`]

	// only allow a photo if we also have a photo_thumb
	if havePhoto != havePhotoThumb {
		delete(properties, `photo`)
		delete(properties, `photo_thumb`)
	}

	// loop through allowed properties
	allowedProperties := []string{"photo", "photo_thumb", "allow_email", "allow_post"}
	for _, key := range allowedProperties {
		if _, ok := properties[key]; ok {
			filteredProperties[key] = properties[key]
		}
	}
	return filteredProperties
}

// PUT - Wrapper for UserProperties::SetUserProperty
func putGameUser(r *http.Request) (user *User, appErr *ApplicationError) {
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
	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id " + vars["game_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	decoder := json.NewDecoder(r.Body)
	var userPost UserPost
	err := decoder.Decode(&userPost)
	if err != nil {
		return nil, NewApplicationError("Invalid JSON", err, ErrCodeInvalidJSON)
	}

	user, appErr = GetUserById(userId)
	if appErr != nil {
		return nil, appErr
	}

	filteredProperties := filterProperties(userPost.Properties)
	appErr = user.SetUserProperties(filteredProperties)
	if appErr != nil {
		return nil, appErr
	}

	// Check if we have an email to change
	email := userPost.Email
	if email == "" {
		return nil, nil
	}

	// Change the user's email
	appErr = user.ChangeEmail(email)
	if appErr != nil {
		return nil, appErr
	}

	return nil, nil
}

type GameUserPost struct {
	GamePassword string `json:"game_password"`
	TeamId       string `json:"team_id"`
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
		msg := "Invalid UUID: user_id " + vars["user_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}
	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id " + vars["game_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	user, appErr := GetUserById(userId)
	if appErr != nil {
		return nil, appErr
	}

	decoder := json.NewDecoder(r.Body)
	var gameUserPost GameUserPost
	err := decoder.Decode(&gameUserPost)
	if err != nil {
		return nil, NewApplicationError("Invalid JSON", err, ErrCodeInvalidJSON)
	}

	gameMapping, appErr = user.JoinGame(gameId, gameUserPost.GamePassword)
	if appErr != nil {
		return nil, appErr
	}

	teamId := uuid.Parse(gameUserPost.TeamId)
	if teamId == nil {
		return gameMapping, nil
	}

	// if there was an error joining a team, fail silently
	sentryRequest := raven.NewHttp(r)
	user, appErr = GetUserById(userId)
	if appErr != nil {
		LogWithSentry(appErr, nil, raven.WARNING, map[string]interface{}{"user_id": userId.String(), "game_id": gameId.String(), "team_id": teamId.String()}, sentryRequest)
		return gameMapping, nil
	}

	sentryUser := NewSentryUser(user)
	gameMapping, appErr = user.JoinTeam(teamId)
	if appErr != nil {
		LogWithSentry(appErr, nil, raven.WARNING, map[string]interface{}{"game_id": gameId.String(), "team_id": teamId.String()}, sentryUser, sentryRequest)
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
		msg := "Invalid UUID: user_id " + vars["user_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id " + vars["game_id"]
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
