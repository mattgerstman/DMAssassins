package main

import (
	"net/http"
)
//Consult the UserHandler for how I'm actually handling Handlers right now
func GameHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		//case "GET":
		//	WriteObjToPayload(w, getUser(w, r))
		case "POST":
			obj, err := assignTargets()
		// case "DELETE":
		// 	WriteObjToPayload(w, deleteUser(w, r))
		default:
			obj = nil
			err = NewSimpleApplicationError("Invalid Http Method", ERROR_INVALID_METHOD)
		}
		WriteObjToPayload(w, obj, err)
	}
}
