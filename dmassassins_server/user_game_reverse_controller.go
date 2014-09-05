package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

func getNewUserGame(r *http.Request) ([]*Game, *ApplicationError) {
	appErr := RequiresLogin(r)
	if appErr != nil {
		return nil, appErr
	}
	vars := mux.Vars(r)
	userId := uuid.Parse(vars["user_id"])

	user, appErr := GetUserById(userId)
	if appErr != nil {
		return nil, appErr
	}

	return user.GetNewGamesForUser()
}

// Handler for /game path
func UserGameNewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "GET":
			obj, err = getNewUserGame(r)
		default:
			obj = nil
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeInvalidMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
