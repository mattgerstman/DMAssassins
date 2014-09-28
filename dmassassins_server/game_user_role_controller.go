package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

// POST - Wrapper for GameMapping:ChangeRole
func postGameUserRole(r *http.Request) (gameMapping *GameMapping, appErr *ApplicationError) {
	appErr = RequiresAdmin(r)
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

	gameMapping, appErr = GetGameMapping(userId, gameId)
	if appErr != nil {
		return nil, appErr
	}

	r.ParseForm()
	role := r.FormValue("role")
	if role == "" {
		msg := "Missing Parameter: role"
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}

	if role == "dm_super_admin" {
		appErr = RequiresSuperAdmin(r)
		if appErr != nil {
			return nil, appErr
		}
	}

	appErr = gameMapping.ChangeRole(role)
	if appErr != nil {
		return nil, appErr
	}

	return gameMapping, nil
}

// Handler for /game path
func GameUserRoleHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "POST":
			obj, err = postGameUserRole(r)

		default:
			obj = nil
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeNotFoundMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
