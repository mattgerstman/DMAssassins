package main

import (
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
	r.ParseForm()
	vars := mux.Vars(r)
	username := vars["username"]
	user, err := GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	return user.GetTarget()
}

// Kill a target, delete User may eventually be used by an admin
func deleteTarget(r *http.Request) (string, *ApplicationError) {

	fmt.Println(r)
	vars := mux.Vars(r)
	username := vars["username"]

	r.ParseForm()
	secret := r.Header.Get("X-DMAssassins-Secret")

	if secret == "" {
		msg := "Missing Parameter: secret."
		err := errors.New("Missing Parameter")
		return "", NewApplicationError(msg, err, ErrCodeMissingParameter)
	}

	fmt.Println(secret)
	//need to actually handle the case where the user doesn't exist
	user, err := GetUserByUsername(username)
	if err != nil {
		return "", err
	}
	return user.KillTarget(secret)
}

// Assigns targets, needs to be updated to only allow admins
func postTarget(r *http.Request) (map[string]string, *ApplicationError) {

	return AssignTargets()

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
		case "POST":
			obj, err = postTarget(r)
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
