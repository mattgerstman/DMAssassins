package main

import (
	"errors"
	"net/http"
)

func postGame() (interface{}, *ApplicationError) {
	AssignTargets()
	return nil, nil
}

//Consult the UserHandler for how I'm actually handling Handlers right now
func GameHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		//case "GET":
		//	WriteObjToPayload(w, getUser(w, r))
		case "POST":
			obj, err = postGame()
		// case "DELETE":
		// 	WriteObjToPayload(w, deleteUser(w, r))
		default:
			obj = nil
			msg := "Not Found"
			err := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, err, ErrCodeInvalidMethod)
		}
		WriteObjToPayload(w, obj, err)
	}
}
