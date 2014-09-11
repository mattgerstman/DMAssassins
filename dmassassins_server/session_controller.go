package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"net/http"
)

// POST - Takes data from facebook and returns an authenticated user/game
func postSession(w http.ResponseWriter, r *http.Request) (response map[string]interface{}, appErr *ApplicationError) {

	// Parse facebook id and token from form
	r.ParseForm()
	facebookId := r.FormValue("facebook_id")
	if facebookId == "" {
		msg := "Missing Parameter: facebook_id."
		err := errors.New("Missing Parameter")
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}
	facebookToken := r.FormValue("facebook_token")
	if facebookToken == "" {
		msg := "Missing Parameter: facebook_token."
		err := errors.New("Missing Parameter")
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}

	// Get the user data from the facebook data
	user, appErr := GetUserFromFacebookData(facebookId, facebookToken)
	if appErr != nil {
		return nil, appErr
	}

	// Start building out the response
	response = make(map[string]interface{})
	response["user"] = user

	// Get the current db token to pass down to the user
	token, appErr := user.GetToken()
	if appErr != nil {
		return nil, appErr
	}
	response["token"] = token

	response["game_mapping"] = nil
	response["game"] = nil

	// If we have a gameId try to get the game mapping first from that
	gameId := uuid.Parse(r.FormValue("game_id"))
	var gameMapping *GameMapping
	if gameId != nil {
		gameMapping, appErr = GetGameMapping(user.UserId, gameId)
		if appErr != nil {
			if appErr.Code != ErrCodeNotFoundGameMapping {
				return nil, appErr
			}
			return response, nil
		}
	}

	// if we don't have an appropriate game mapping get an arbirtary one
	if gameMapping == nil {
		gameMapping, appErr = user.GetArbitraryGameMapping()
		if appErr != nil {
			if appErr.Code != ErrCodeNoGameMappings {
				return nil, appErr
			}
			// If we aren't mapped to any game return an appropriate response
			return response, nil
		}
	}

	response["game_mapping"] = gameMapping
	// Get the game for whatever game mapping we're using
	game, appErr := GetGameById(gameMapping.GameId)
	if appErr != nil {
		return nil, appErr
	}

	response["game"] = game
	return response, nil
}

func SessionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "POST":
			obj, err = postSession(w, r)
		default:
			obj = nil
			msg := "Not Found"
			err := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, err, ErrCodeNotFoundMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
