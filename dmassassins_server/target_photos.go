package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
)

// GET function for /user/{user_id}/target/photos returns a user's information
func getTargetPhotos(r *http.Request) (photos []interface{}, appErr *ApplicationError) {
	_, appErr = RequiresUser(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)
	userId := uuid.Parse(vars["user_id"])
	if userId == nil {
		msg := "Invalid UUID: user_id " + vars["user_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}
	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id " + vars["game_id"]
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	// get user
	user, err := GetUserById(userId)
	if err != nil {
		return nil, err
	}

	// get target
	target, appErr := user.GetTarget(gameId)
	if appErr != nil {
		return nil, appErr
	}

	return target.GetFacebookPhotos()
}

// Handler for /user/{user_id}/target
func TargetPhotosHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "GET":
			obj, err = getTargetPhotos(r)
		default:
			obj = nil
			msg := "Not Found"
			err := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, err, ErrCodeNotFoundMethod)

		}
		WriteObjToPayload(w, r, obj, err)
	}
}
