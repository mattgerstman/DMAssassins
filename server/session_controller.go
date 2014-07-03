package main

import (
	_ "github.com/lib/pq"
	"net/http"
	//"fmt"
	//"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("some-thing-very-secret"))

func getSession(w http.ResponseWriter, r *http.Request) interface{} {
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")
	user := GetUserByEmail(email)
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

	return valid
}

func killSession(w http.ResponseWriter, r *http.Request) interface{} {
	session, _ := store.Get(r, "DMAssassins")
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	return session.Save(r, w)
}

func SessionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "GET":
			WriteObjToPayload(w, getSession(w, r))
			//case "POST":
		//WriteObjToPayload(w, postUser(w, r))

		//servePostUser(db)(w, r)
		case "DELETE":
			WriteObjToPayload(w, killSession(w, r))
		default:
		}
	}
}
