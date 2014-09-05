package main

import (
	"errors"
	"fmt"
	//"github.com/gorilla/sessions"
	"net/http"
)

// Takes data from facebook and returns an authenticated user
func postSession(w http.ResponseWriter, r *http.Request) (interface{}, *ApplicationError) {
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

	user, appErr := GetUserFromFacebookData(facebookId, facebookToken)
	if appErr != nil {
		return nil, appErr
	}

	response := make(map[string]interface{})

	response["user"] = user
	token, appErr := user.GetToken()
	if appErr != nil {
		return nil, appErr
	}
	response["token"] = token

	game, appErr := user.GetArbitraryGame()
	if (appErr != nil) && (appErr.Code != ErrCodeNoGameMappings) {
		return nil, appErr
	}

	response["game"] = game
	return response, nil
}

func SessionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("SessionHandler()")
		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "POST":
			obj, err = postSession(w, r)
		default:
			obj = nil
			msg := "Not Found"
			err := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, err, ErrCodeInvalidMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
