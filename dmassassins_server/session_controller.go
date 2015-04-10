package main

import (
	"errors"
	"net/http"
)

// POST - Takes data from facebook and returns an authenticated user/game
func postSession(w http.ResponseWriter, r *http.Request) (response map[string]interface{}, appErr *ApplicationError) {
	params, appErr := NewParams(r)
	if appErr != nil {
		return nil, appErr
	}

	// Parse facebook id and token from form
	facebookId, appErr := params.GetStringParam("facebook_id")
	if appErr != nil {
		return nil, appErr
	}

	facebookToken, appErr := params.GetStringParam("facebook_token")
	if appErr != nil {
		return nil, appErr
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

	// Set all of the following to nil if we don't have them yet
	response["game"] = nil

	// If we have a gameId try to get the game mapping first from that
	gameId, _ := params.GetUUIDParam("game_id")
	var gameMapping *GameMapping
	if gameId != nil {
		gameMapping, appErr = GetGameMapping(user.UserId, gameId)
		if appErr != nil {
			// if we have no games return here
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
			// if we have no games return here
			if appErr.Code != ErrCodeNoGameMappings {
				return nil, appErr
			}
			// If we aren't mapped to any game return an appropriate response
			return response, nil
		}
	}

	// Get the game for whatever game mapping we're using
	game, appErr := GetGameById(gameMapping.GameId)
	if appErr != nil {
		return nil, appErr
	}

	response["game"] = game

	return response, nil
}

// Handler for the session controller
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
