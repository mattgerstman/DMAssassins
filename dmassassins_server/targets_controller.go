package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	//"github.com/gorilla/schema"
)

// GET function for /users/{username}/target returns a user's information
// Need to add permissions to this on a per user basis
func getTarget(r *http.Request) (*User, *ApplicationError) {
	appErr := RequiresUser(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)
	userId := uuid.Parse(vars["user_id"])
	user, err := GetUserById(userId)
	if err != nil {
		return nil, err
	}
	gameId := uuid.Parse(vars["game_id"])
	return user.GetTarget(gameId)
}

// Kill a target, delete User may eventually be used by an admin
func deleteTarget(r *http.Request) (uuid.UUID, *ApplicationError) {
	vars := mux.Vars(r)

	userId := uuid.Parse(vars["user_id"])
	appErr := RequiresUser(r)
	if appErr != nil {
		return nil, appErr
	}

	r.ParseForm()
	secret := r.Header.Get("X-DMAssassins-Secret")

	if secret == "" {
		msg := "Missing Header: X-DMAssassins-Secret."
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeMissingHeader)
	}

	user, err := GetUserById(userId)
	if err != nil {
		return nil, err
	}
	gameId := uuid.Parse(vars["game_id"])
	return user.KillTarget(gameId, secret)
}

// Handler for /user/{username}/target
func TargetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("TargetHandler()")
		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "GET":
			obj, err = getTarget(r)
		//case "POST":
		//obj, err = postTarget(r)
		case "DELETE":
			obj, err = deleteTarget(r)
		default:
			obj = nil
			msg := "Not Found"
			err := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, err, ErrCodeInvalidMethod)

		}
		WriteObjToPayload(w, r, obj, err)
	}
}
