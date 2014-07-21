package main

import (
	"net/http"
)

func GameHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		//case "GET":
		//	WriteObjToPayload(w, getUser(w, r))
		case "POST":
			obj, err := assignTargets()
			WriteObjToPayload(w, obj, err)
		// case "DELETE":
		// 	WriteObjToPayload(w, deleteUser(w, r))
		default:
			obj = nil
			err = NewSimpleApplicationError("Invalid Http Method", ERROR_INVALID_METHOD)
		}
		WriteObjToPayload(w, obj, err)
	}
}
