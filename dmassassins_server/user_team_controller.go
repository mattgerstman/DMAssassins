package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

func postUserTeam(r *http.Request) (*GameMapping, *ApplicationError) {
	appErr := RequiresAdmin(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)
	userId := uuid.Parse(vars["user_id"])

	user, appErr := GetUserById(userId)
	if appErr != nil {
		return nil, appErr
	}

	teamId := uuid.Parse(vars["team_id"])
	return user.JoinTeam(teamId)
}

// func getUserTeam(r *http.Request) (*Team, *ApplicationError) {
// 	appErr := RequiresCaptain(r)
// 	if appErr != nil {
// 		return nil, appErr
// 	}

// 	vars := mux.Vars(r)
// 	teamId := uuid.Parse(vars["team_id"])
// 	return GetTeamById(teamId)
// }

func deleteUserTeam(r *http.Request) (*GameMapping, *ApplicationError) {
	appErr := RequiresAdmin(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)

	userId := uuid.Parse(vars["user_id"])
	user, appErr := GetUserById(userId)
	if appErr != nil {
		return nil, appErr
	}

	teamId := uuid.Parse(vars["team_id"])
	return user.LeaveTeam(teamId)
}

// Handler for /team path
func UserTeamHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		//case "GET":
			//obj, err = getTeamId(r)
		case "POST":
			obj, err = postTeamId(r)
		case "DELETE":
			obj, err = nil, deleteTeamId(r)
		default:
			obj = nil
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeInvalidMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
