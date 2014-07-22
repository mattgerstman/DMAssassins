package main

import (
	_ "github.com/lib/pq"
	"net/http"
	//"fmt"
	//"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
)

//yes i know i need a real secret key and i should read it from a config file
var store = sessions.NewCookieStore([]byte("some-thing-very-secret"))

//Tell me what's wrong with this
func getSession(w http.ResponseWriter, r *http.Request) (interface{}, *ApplicationError) {
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
		case "GET":
			obj, err = getSession(w, r)
			//case "POST":
		//WriteObjToPayload(w, postUser(w, r))

		//servePostUser(db)(w, r)
		case "DELETE":
			obj, err = killSession(w, r)
		default:
			obj = nil
			err = NewSimpleApplicationError("Invalid Http Method", ERROR_INVALID_METHOD)
		}
		WriteObjToPayload(w, obj, err)
	}
}
