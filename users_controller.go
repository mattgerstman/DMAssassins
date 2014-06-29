package main

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	//"github.com/gorilla/schema"
)

func getUser(w http.ResponseWriter, r *http.Request) *User {
	vars := mux.Vars(r)
	email := vars["email"]

	if email == "" {
		http.Error(w, "Missing Parameter: email.", http.StatusBadRequest)
		return nil
	}

	return GetUserByEmail(email)
}

func postUser(w http.ResponseWriter, r *http.Request) *User {
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
	errorMessage := fmt.Sprintf("Missing Parameter: %s", missingParam)
	if missingParam != "" {
		http.Error(w, errorMessage, http.StatusBadRequest)
		return nil
	}
	return NewUser(email, password, secret)
}

func deleteUser(w http.ResponseWriter, r *http.Request) string {
	session, _ := store.Get(r, "DMAssassins")
	logged_in_user, ok := session.Values["user_id"].(string)

	if !ok || logged_in_user == "" {
		errorMessage := "Error: Not logged in"
		http.Error(w, errorMessage, http.StatusBadRequest)
		return ""
	}

	r.ParseForm()
	secret := r.FormValue("secret")

	user := GetUserById(logged_in_user)
	return user.KillTarget(secret)
}

func UserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "GET":
			WriteObjToPayload(w, getUser(w, r))
		case "POST":
			WriteObjToPayload(w, postUser(w, r))
		case "DELETE":
			WriteObjToPayload(w, deleteUser(w, r))
		default:
		}
	}
}
