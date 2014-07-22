package main

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	//"github.com/gorilla/schema"
)

func getUser(w http.ResponseWriter, r *http.Request) (*User, *ApplicationError) {
	vars := mux.Vars(r)
	email := vars["email"]

	if email == "" {
		msg := "Missing Parameter: email."		
		return nil, NewSimpleApplicationError(msg, ERROR_MISSING_PARAMETER)
	}

	return GetUserByEmail(email)
}

func postUser(w http.ResponseWriter, r *http.Request) (*User, *ApplicationError) {
	r.ParseForm()
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	secret := r.PostFormValue("secret")

	missingParam := ""
	switch {
	case email == "":
		missingParam = "email"
	case password == "":
		missingParam = "password"
	case secret == "":
		missingParam = "secret"
	}
	msg := fmt.Sprintf("Missing Parameter: %s", missingParam)
	if missingParam != "" {
		return nil, NewSimpleApplicationError(msg, ERROR_MISSING_PARAMETER)
	}
	return NewUser(email, password, secret)
}

func deleteUser(w http.ResponseWriter, r *http.Request) (string, *ApplicationError) {
	session, _ := store.Get(r, "DMAssassins")
	logged_in_user, ok := session.Values["user_id"].(string)

	if !ok || logged_in_user == "" {
		msg := "Error: Not logged in"
		return "", NewSimpleApplicationError(msg, ERROR_NO_SESSION)
	}

	r.ParseForm()
	secret := r.FormValue("secret")

	user, err := GetUserById(logged_in_user)
	_ = err
	return user.KillTarget(secret)
}

func UserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "GET":
			obj, err = getUser(w, r)
		case "POST":
			obj, err = postUser(w, r)
		case "DELETE":
			obj, err = deleteUser(w, r)
		default:
			obj = nil
			err = NewSimpleApplicationError("Invalid Http Method", ERROR_INVALID_METHOD)

		}
		WriteObjToPayload(w, obj, err)
	}
}
