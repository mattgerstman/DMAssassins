package main

import (
	"errors"
	"fmt"
	//"github.com/gorilla/sessions"
	"net/http"
)

//yes i know i need a real secret key and i should read it from a config file
// var store = sessions.NewCookieStore([]byte("some-thing-very-secret"))

// Takes data from facebook and returns an authenticated user
func postSession(w http.ResponseWriter, r *http.Request) (interface{}, *ApplicationError) {
	r.ParseForm()
	facebook_id := r.FormValue("facebook_id")
	facebook_token := r.FormValue("facebook_token")
	user, err := GetUserFromFacebookData(facebook_id, facebook_token)

	return user, err
}

// // Kill a session this will probably be rewritten later with basic auth
// func deleteSession(w http.ResponseWriter, r *http.Request) (interface{}, *ApplicationError) {
// 	session, _ := store.Get(r, "DMAssassins")
// 	session.Options = &sessions.Options{
// 		Path:     "/",
// 		MaxAge:   -1,
// 		HttpOnly: true,
// 	}

// 	return session.Save(r, w), nil
// }

func SessionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("SessionHandler()")
		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "POST":
			obj, err = postSession(w, r)
		// case "DELETE":
		// 	obj, err = deleteSession(w, r)
		default:
			obj = nil
			msg := "Not Found"
			err := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, err, ErrCodeInvalidMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
