package main

import (
	"net/http"
)

func GameHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		//case "GET":
		//	WriteObjToPayload(w, getUser(w, r))
		case "POST":
			WriteObjToPayload(w, assignTargets())
		// case "DELETE":
		// 	WriteObjToPayload(w, deleteUser(w, r))
		default:
		}
	}
}
