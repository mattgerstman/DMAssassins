package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	//"github.com/gorilla/schema"
)

// GET function for /users/{username} returns a user's information
// Need to add permissions to this on a per user basis
func getUser(r *http.Request) (*User, *ApplicationError) {
	r.ParseForm()
	vars := mux.Vars(r)
	username := vars["username"]
	fmt.Println(username)

	return GetUserByUsername(username)
}

// Create user, need to create an auth token system for signups
func postUser(r *http.Request) (*User, *ApplicationError) {
	r.ParseForm()
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	email := r.PostFormValue("email")
	secret := r.PostFormValue("secret")

	missingParam := ""
	switch {
	case username == "":
		missingParam = "username"
	case password == "":
		missingParam = "password"
	case secret == "":
		missingParam = "secret"
	}
	msg := fmt.Sprintf("Missing Parameter: %s", missingParam)
	err := errors.New("Missing Parameter")
	if missingParam != "" {
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}
	return NewUser(username, email, password, secret)
}

// Kill a target, delete User may eventually be used by an admin
func deleteUser(r *http.Request) (string, *ApplicationError) {
	session, _ := store.Get(r, "DMAssassins")
	logged_in_user, ok := session.Values["user_id"].(string)

	if !ok || logged_in_user == "" {
		msg := "Error: Not logged in"
		err := errors.New("No session found for user")
		return "", NewApplicationError(msg, err, ErrCodeNoSession)
	}

	r.ParseForm()
	secret := r.FormValue("secret")

	fmt.Println(secret)
	//need to actually handle the case where the user doesn't exist
	user, err := GetUserById(logged_in_user)
	_ = err
	return user.KillTarget(secret)
}

// Handler for /users/ and /users/{username}/
func UserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("UserHandler()")
		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "GET":
			obj, err = getUser(r)
		case "POST":
			obj, err = postUser(r)
		case "DELETE":
			obj, err = deleteUser(r)
		default:
			obj = nil
			msg := "Not Found"
			err := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, err, ErrCodeInvalidMethod)

		}
		WriteObjToPayload(w, r, obj, err)
	}
}
