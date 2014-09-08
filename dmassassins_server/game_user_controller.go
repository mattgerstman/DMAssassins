package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
)

// POST - Wrapper for GameMapping:JoinGame
func postGameUser(r *http.Request) (gameMapping *GameMapping, appErr *ApplicationError) {
	appErr = RequiresLogin(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)
	userId := uuid.Parse(vars["user_id"])
	if userId == nil {
		msg := "Invalid UUID: user_id" + userId.String()
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}
	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id" + gameId.String()
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}
	gamePassword := r.FormValue("game_password")

	user, appErr := GetUserById(userId)
	if appErr != nil {
		return nil, appErr
	}

	return user.JoinGame(gameId, gamePassword)
}

// GET - Wrapper for GameMapping:GetGameMapping, usually used for user_role or alive/kill status
func getGameUser(r *http.Request) (user *User, appErr *ApplicationError) {
	//appErr = RequiresUser(r)
	if appErr != nil {
		return nil, appErr
	}
	vars := mux.Vars(r)
	userId := uuid.Parse(vars["user_id"])
	if userId == nil {
		msg := "Invalid UUID: user_id" + userId.String()
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id" + gameId.String()
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	user, appErr = GetUserById(userId)
	if appErr != nil {
		return nil, appErr
	}

	gameMapping, appErr := GetGameMapping(userId, gameId)
	if appErr != nil {
		return nil, appErr
	}

	user.Properties["secret"] = gameMapping.Secret
	user.Properties["team"] = ""

	fmt.Println(gameMapping.TeamId.String());

	if gameMapping.TeamId == nil {
		return user, nil
	}

	team, appErr := GetTeamById(gameMapping.TeamId)
	if appErr != nil {
		return nil, appErr
	}
	user.Properties["team"]= team.TeamName
	
	return user, nil
}

// Handler for /game path
func GameUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "GET":
			obj, err = getGameUser(r)

		case "POST":
			obj, err = postGameUser(r)
		default:
			obj = nil
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeNotFoundMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
