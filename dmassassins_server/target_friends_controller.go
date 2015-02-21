package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
)

// GET function for /user/{user_id}/target/friends returns a user's information
func getTargetFriends(r *http.Request) (friendData map[string]interface{}, appErr *ApplicationError) {
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

	friendData = make(map[string]interface{})

	// query db for mutual friends
	friends, count, appErr := user.GetMutualFriends(target.FacebookId)
	if appErr != nil {
		return nil, appErr
	}

	friendData[`friends`] = friends
	friendData[`count`] = count

	// if we have friends return them
	if count != 0 {
		// query facebook for friends for future requests
		go user.StoreUserFriends()
		return friendData, nil
	}

	// if we have no friends query facebook and try again
	user.StoreUserFriends()
	target.StoreUserFriends()

	// query db for mutual friends
	friendData[`friends`], friendData[`count`], appErr = user.GetMutualFriends(target.FacebookId)
	if appErr != nil {
		return nil, appErr
	}

	return friendData, nil
}

// Handler for /user/{user_id}/target
func TargetFriendsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "GET":
			obj, err = getTargetFriends(r)
		default:
			obj = nil
			msg := "Not Found"
			err := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, err, ErrCodeNotFoundMethod)

		}
		WriteObjToPayload(w, r, obj, err)
	}
}
