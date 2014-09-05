package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

func postTeamId(r *http.Request) (*Team, *ApplicationError) {
	appErr := RequiresAdmin(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)
	teamId := uuid.Parse(vars["team_id"])

	team, appErr := GetTeamById(teamId)
	if appErr != nil {
		return nil, appErr
	}

	teamName := r.FormValue("team_name")
	return team.Rename(teamName)
}

func getTeamId(r *http.Request) (*Team, *ApplicationError) {
	appErr := RequiresCaptain(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)
	teamId := uuid.Parse(vars["team_id"])
	return GetTeamById(teamId)
}

func deleteTeamId(r *http.Request) (*ApplicationError) {
	appErr := RequiresAdmin(r)
	if appErr != nil {
		return appErr
	}

	vars := mux.Vars(r)
	teamId := uuid.Parse(vars["team_id"])
	return DeleteTeam(teamId)
}

// Handler for /team path
func TeamIdHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "GET":
			obj, err = getTeamId(r)
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
