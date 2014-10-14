package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type GameSettingsPut struct {
	GameName     string `json:"game_name"`
	GamePassword string `json:"game_password"`
	EnableTeams  string `json:"game_teams_enabled"`
}

// PUT - Changes game settings
func putGameId(r *http.Request) (game *Game, appErr *ApplicationError) {
	_, appErr = RequiresAdmin(r)
	if appErr != nil {
		return nil, appErr
	}

	// Get Game Id
	vars := mux.Vars(r)
	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id " + vars["game_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	game, appErr = GetGameById(gameId)
	if appErr != nil {
		return nil, appErr
	}

	decoder := json.NewDecoder(r.Body)
	var gameSettingsPut GameSettingsPut
	err := decoder.Decode(&gameSettingsPut)
	if err != nil {
		return nil, NewApplicationError("Invalid JSON", err, ErrCodeInvalidJSON)
	}

	// Rename Game
	gameName := gameSettingsPut.GameName
	if gameName != "" {
		game.Rename(gameName)
	}

	// Change password
	game.ChangePassword(gameSettingsPut.GamePassword)

	//Set teams enabled
	game.SetGameProperty("teams_enabled", gameSettingsPut.EnableTeams)

	return game, nil
}

// POST - Starts a game
func postGameId(r *http.Request) (game *Game, appErr *ApplicationError) {
	_, appErr = RequiresAdmin(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)
	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id" + vars["game_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	game, appErr = GetGameById(gameId)
	if appErr != nil {
		return nil, appErr
	}

	appErr = game.Start()
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = game.sendStartGameEmail()
	if appErr != nil {
		return nil, appErr
	}

	return game, nil
}

// GET - Gets a game
func getGameId(r *http.Request) (game *Game, appErr *ApplicationError) {
	role, appErr := RequiresUser(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)
	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id" + vars["game_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	game, appErr = GetGameById(gameId)

	if appErr != nil {
		return nil, appErr
	}

	if !CompareRole(role, RoleAdmin) {
		return game, nil
	}

	password, appErr := game.GetPassword()
	if appErr != nil {
		return nil, appErr
	}

	game.Properties["game_password"] = password
	return game, nil

}

// DELETE - Ends a game
func deleteGameId(r *http.Request) (game *Game, appErr *ApplicationError) {
	_, appErr = RequiresAdmin(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)
	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id" + vars["game_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	game, appErr = GetGameById(gameId)

	if appErr != nil {
		return nil, appErr
	}

	appErr = game.End()
	if appErr != nil {
		return nil, appErr
	}
	return game, nil

}

// Handler for /game path
func GameIdHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "GET":
			obj, err = getGameId(r)
		case "POST":
			obj, err = postGameId(r)
		case "PUT":
			obj, err = putGameId(r)
		case "DELETE":
			obj, err = deleteGameId(r)
		default:
			obj = nil
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeNotFoundMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
