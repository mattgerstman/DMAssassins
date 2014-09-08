package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
)

// GET function for /users/{username}/target returns a user's information
// Need to add permissions to this on a per user basis
func getTarget(r *http.Request) (user *User, appErr *ApplicationError) {
	appErr = RequiresUser(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)
	userId := uuid.Parse(vars["user_id"])
	if userId == nil {
		msg := "Invalid UUID: user_id" + userId.String()
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}
	user, err := GetUserById(userId)
	if err != nil {
		return nil, err
	}
	gameId := uuid.Parse(vars["game_id"])
	if gameId == nil {
		msg := "Invalid UUID: game_id" + gameId.String()
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}


	user, appErr = user.GetTarget(gameId)
	if appErr != nil {
		return nil, appErr
	}


	gameMapping, appErr := GetGameMapping(userId, gameId)
	if appErr != nil {
		return nil, appErr
	}
	user.Properties["team"] = ""

	fmt.Println(gameMapping.TeamId.String());

	if gameMapping.TeamId == nil {
		return user, nil
	}

	team, appErr := GetTeamById(gameMapping.TeamId)
	if appErr != nil {
		return nil, appErr
	}
	user.Properties["team"]= team.TeamName

	return user, nil
}

// DELETE - Kill a target, delete User may eventually be used by an admin
func deleteTarget(r *http.Request) (targetId uuid.UUID, appErr *ApplicationError) {
	appErr = RequiresUser(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)
	userId := uuid.Parse(vars["user_id"])
	if userId == nil {
		msg := "Invalid UUID: user_id" + userId.String()
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	r.ParseForm()
	secret := r.Header.Get("X-DMAssassins-Secret")

	if secret == "" {
		msg := "Missing Header: X-DMAssassins-Secret."
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeMissingHeader)
	}

	user, err := GetUserById(userId)
	if err != nil {
		return nil, err
	}
	gameId := uuid.Parse(vars["game_id"])
	return user.KillTarget(gameId, secret, true)
}

// Handler for /user/{username}/target
func TargetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("TargetHandler()")
		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "GET":
			obj, err = getTarget(r)
		//case "POST":
		//obj, err = postTarget(r)
		case "DELETE":
			obj, err = deleteTarget(r)
		default:
			obj = nil
			msg := "Not Found"
			err := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, err, ErrCodeNotFoundMethod)

		}
		WriteObjToPayload(w, r, obj, err)
	}
}
