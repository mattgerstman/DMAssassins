package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	//"github.com/gorilla/schema"
)

func getUser(r *http.Request) (*User, *ApplicationError) {
	r.ParseForm()
	vars := mux.Vars(r)
	email := vars["email"]

	if email == "" {
		msg := "Missing Parameter: email."
		err := errors.New("Missing Parameter")
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}

	return GetUserByEmail(email)
}

func postUser(r *http.Request) (*User, *ApplicationError) {
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
	err := errors.New("Missing Parameter")
	if missingParam != "" {
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}
	return NewUser(email, password, secret)
}

//None of my functions need a ResponseWriter anymore but I haven't removed it yet
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

//This is pretty much the default of how I'm writing my Handlers. I'm unaware of anything wrong with it.
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
