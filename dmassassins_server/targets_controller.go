package main

import (
	"errors"
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	//"github.com/gorilla/schema"
)

// GET function for /users/{username}/target returns a user's information
// Need to add permissions to this on a per user basis
func getTarget(r *http.Request) (*User, *ApplicationError) {
	r.ParseForm()
	vars := mux.Vars(r)
	userId := uuid.Parse(vars["user_id"])
	user, err := GetUserById(userId)
	if err != nil {
		return nil, err
	}

	return user.GetTarget()
}

// Kill a target, delete User may eventually be used by an admin
func deleteTarget(r *http.Request) (uuid.UUID, *ApplicationError) {

	fmt.Println(r)
	vars := mux.Vars(r)
	
	userId := uuid.Parse(vars["user_id"])

	r.ParseForm()
	secret := r.Header.Get("X-DMAssassins-Secret")

	if secret == "" {
		msg := "Missing Header: X-DMAssassins-Secret."
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeMissingHeader)
	}

	gameId := uuid.Parse(vars["game_id"])

	user, err := GetUserById(userId)
	if err != nil {
		return nil, err
	}
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
