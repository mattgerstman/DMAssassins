package main

import (
	//"errors"
	"net/http"
)

// Handler for /game path
func GameHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// This handler will be used when I implement plot twists and multiple games
		// Fuck You Taylor

		// var obj interface{}
		// var err *ApplicationError

		// switch r.Method {
		// case "POST":
		// 	obj, err = postGame()
		// default:
		// 	obj = nil
		// 	msg := "Not Found"
		// 	tempErr := errors.New("Invalid Http Method")
		// 	err = NewApplicationError(msg, tempErr, ErrCodeInvalidMethod)
		// }
		// WriteObjToPayload(w, r, obj, err)
	}
}
