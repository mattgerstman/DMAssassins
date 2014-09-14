package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type NewGamePost struct {
	GameName     string `json:"game_name"`
	GamePassword string `json:"game_id"`
}

// POST - Controller Wrapper for Game:NewGame
func postUserGame(r *http.Request) (game *Game, appErr *ApplicationError) {
	appErr = RequiresLogin(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)
	userId := uuid.Parse(vars["user_id"])
	if userId == nil {
		msg := "Invalid Parameter: user_id " + vars["user_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
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
	return NewGame(gameName, userId, gamePassword)
}

// GET - gets a list of games a user is in
func getUserGame(r *http.Request) (response map[string][]*Game, appErr *ApplicationError) {
	appErr = RequiresLogin(r)
	if appErr != nil {
		return nil, appErr
	}
	vars := mux.Vars(r)
	userId := uuid.Parse(vars["user_id"])
	if userId == nil {
		msg := "Invalid UUID: user_id" + userId.String()
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}

	user, appErr := GetUserById(userId)
	if appErr != nil {
		return nil, appErr
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
func UserGameHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "POST":
			obj, err = postUserGame(r)
		case "PUT":
			obj, err = postUserGame(r)
		case "GET":
			obj, err = getUserGame(r)
		default:
			obj = nil
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeNotFoundMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
