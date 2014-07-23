package main

import (
	"errors"
	"fmt"
	"net/http"
	//"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
)

//yes i know i need a real secret key and i should read it from a config file
var store = sessions.NewCookieStore([]byte("some-thing-very-secret"))

//Tell me what's wrong with this
func postSession(w http.ResponseWriter, r *http.Request) (interface{}, *ApplicationError) {
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := GetUserByEmail(email)

	if err != nil {
		return nil, err
	}

	valid := user.CheckPassword(password)

	if valid {
		session, _ := store.Get(r, "DMAssassins")
		session.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   1800,
			HttpOnly: true,
		}
		session.Values["user_id"] = user.User_id
		session.Save(r, w)
	}
	fmt.Println(valid)

	return valid, nil
}

//Tell me what's wrong with this
func killSession(w http.ResponseWriter, r *http.Request) (interface{}, *ApplicationError) {
	session, _ := store.Get(r, "DMAssassins")
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	return session.Save(r, w), nil
}

//Consult the UserHandler for how I'm actually handling Handlers right now
func SessionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "POST":
			obj, err = postSession(w, r)
			//case "POST":
		//WriteObjToPayload(w, postUser(w, r))

		//servePostUser(db)(w, r)
		case "DELETE":
			obj, err = killSession(w, r)
		default:
			obj = nil
			msg := "Not Found"
			err := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, err, ErrCodeInvalidMethod)
		}
		WriteObjToPayload(w, obj, err)
	}
}
