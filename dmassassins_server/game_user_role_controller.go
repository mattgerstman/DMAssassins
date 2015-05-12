package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

// PUT - Wrapper for GameMapping:ChangeRole
func putGameUserRole(r *http.Request) (gameMapping *GameMapping, appErr *ApplicationError) {
	_, appErr = RequiresAdmin(r)
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

	params, appErr := NewParams(r)
	if appErr != nil {
		return nil, appErr
	}

	role, appErr := params.GetStringParam("role")
	if appErr != nil {
		return nil, appErr
	}

	if role == "dm_super_admin" {
		_, appErr = RequiresSuperAdmin(r)
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
		case "PUT":
			obj, err = putGameUserRole(r)

		default:
			obj = nil
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeNotFoundMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
