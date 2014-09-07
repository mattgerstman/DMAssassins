package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

// POST - rename a team
func postTeamId(r *http.Request) (team *Team, appErr *ApplicationError) {
	appErr = RequiresAdmin(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)
	teamId := uuid.Parse(vars["team_id"])
	if teamId == nil {
		msg := "Invalid UUID: team_id" + teamId.String()
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	team, appErr = GetTeamById(teamId)
	if appErr != nil {
		return nil, appErr
	}

	teamName := r.FormValue("team_name")
	if teamName == "" {
		msg := "Missing Parameter: team_name"
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}

	appErr = team.Rename(teamName)
	if appErr != nil {
		return nil, appErr
	}
	return team, nil
}

// GET - get a team by its id
func getTeamId(r *http.Request) (team *Team, appErr *ApplicationError) {
	appErr = RequiresCaptain(r)
	if appErr != nil {
		return nil, appErr
	}

	vars := mux.Vars(r)
	teamId := uuid.Parse(vars["team_id"])
	if teamId == nil {
		msg := "Invalid UUID: team_id" + teamId.String()
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}
	return GetTeamById(teamId)
}

// DELETE - deletes a team
func deleteTeamId(r *http.Request) (appErr *ApplicationError) {
	appErr = RequiresAdmin(r)
	if appErr != nil {
		return appErr
	}

	vars := mux.Vars(r)
	teamId := uuid.Parse(vars["team_id"])
	if teamId == nil {
		msg := "Invalid UUID: team_id" + teamId.String()
		err := errors.New(msg)
		return NewApplicationError(msg, err, ErrCodeMissingParameter)
	}
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
			err = NewApplicationError(msg, tempErr, ErrCodeNotFoundMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
