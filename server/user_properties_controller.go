package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	//"github.com/gorilla/schema"
)

func getUserProperty(r *http.Request) (string, *ApplicationError) {
	r.ParseForm()
	vars := mux.Vars(r)
	username := vars["username"]
	key := vars["key"]

	if username == "" {
		msg := "Missing Parameter: username."
		err := errors.New("Missing Parameter")
		return "", NewApplicationError(msg, err, ErrCodeMissingParameter)
	}
	user, err := GetUserByUsername(username)
	if err != nil {
		return "", err
	}
	return user.GetUserProperty(key)
}

func postUserProperty(r *http.Request) (bool, *ApplicationError) {
	r.ParseForm()
	vars := mux.Vars(r)
	username := vars["username"]
	key := vars["key"]
	value := r.PostFormValue("value")

	if value == "" {
		msg := "Missing Parameter: value."
		err := errors.New("Missing Parameter")
		return false, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}

	user, err := GetUserByUsername(username)
	if err != nil {
		return false, err
	}
	return user.SetUserProperty(key, value)
}
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
