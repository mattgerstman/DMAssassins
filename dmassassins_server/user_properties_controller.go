package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	//"github.com/gorilla/schema"
)

// Get a single user property
func getUserProperty(r *http.Request) (string, *ApplicationError) {
	r.ParseForm()
	vars := mux.Vars(r)
	username := vars["username"]
	key := vars["key"]

	user, err := GetUserByUsername(username)
	if err != nil {
		return "", err
	}
	return user.GetUserProperty(key)
}

// Set a single User Property
func postUserProperty(r *http.Request) (*User, *ApplicationError) {
	r.ParseForm()
	vars := mux.Vars(r)
	username := vars["username"]
	key := vars["key"]
	value := r.PostFormValue("value")

	if value == "" {
		msg := "Missing Parameter: value."
		err := errors.New("Missing Parameter")
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}

	user, err := GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	return user.SetUserProperty(key, value)
}

// Handler for User Property
func UserPropertyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("UserPropertyHandler()")
		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "GET":
			obj, err = getUserProperty(r)
		case "POST":
			obj, err = postUserProperty(r)
		default:
			obj = nil
			msg := "Not Found"
			err := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, err, ErrCodeInvalidMethod)

		}
		WriteObjToPayload(w, r, obj, err)
	}
}
