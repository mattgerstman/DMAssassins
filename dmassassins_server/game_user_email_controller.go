package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	"net/http"
)

// POST - Update email settings
func postGameUserEmail(r *http.Request) (appErr *ApplicationError) {
	_, appErr = RequiresUser(r)
	if appErr != nil {
		return appErr
	}
	vars := mux.Vars(r)
	userId := uuid.Parse(vars["user_id"])
	if userId == nil {
		msg := "Invalid UUID: user_id " + vars["user_id"]
		err := errors.New(msg)
		return NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	user, appErr := GetUserById(userId)
	if appErr != nil {
		return appErr
	}

	// Change email property
	email := r.FormValue("email")
	appErr = user.ChangeEmail(email)
	if appErr != nil {
		return appErr
	}

	// Set allow email
	allowEmail := r.FormValue("allow_email")
	appErr = user.SetUserProperty("allow_email", allowEmail)
	if appErr != nil {
		return appErr
	}

	// Notify the user their email settings have changed
	_, appErr = user.SendChangeEmailEmail()
	if appErr != nil {
		extra := GetExtraDataFromRequest(r)
		LogWithSentry(appErr, map[string]string{"user_id": userId.String()}, raven.WARNING, extra)
	}

	return appErr
}

// Handler for /game path
func GameUserEmailHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "POST":
			err = postGameUserEmail(r)

		default:
			obj = nil
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeNotFoundMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
