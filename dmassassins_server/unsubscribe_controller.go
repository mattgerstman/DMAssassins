package main

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

// GET - Unsubscribes a user from assassins newsletter
func getUnsubscribe(r *http.Request) (msg string, appErr *ApplicationError) {
	vars := mux.Vars(r)
	userId := uuid.Parse(vars["user_id"])
	if userId == nil {
		msg := "Invalid UUID: user_id " + vars["user_id"] + "\n"
		err := errors.New(msg)
		msg += `There was an error unsubscribing, please send an email to `
		msg += `<a href="mailto:` + Config.SupportEmail + `?subject=Unsubscribe">`
		msg += Config.SupportEmail + `</a> to unsubscribe`

		return "", NewApplicationError(msg, err, ErrCodeInvalidUUID)
	}

	user, appErr := GetUserById(userId)
	if appErr != nil {
		appErr.Msg = `There was an error unsubscribing, please send an email to `
		appErr.Msg += `<a href="mailto:` + Config.SupportEmail + `?subject=Unsubscribe">`
		appErr.Msg += Config.SupportEmail + `</a> to unsubscribe`
		return "", appErr
	}

	appErr = user.SetUserProperty("allow_email", "false")
	if appErr != nil {
		appErr.Msg = `There was an error unsubscribing, please send an email to `
		appErr.Msg += `<a href="mailto:` + Config.SupportEmail + `?subject=Unsubscribe: ` + user.Email + `">`
		appErr.Msg += Config.SupportEmail + `</a> to unsubscribe`
		return "", appErr
	}

	success := `Successfully unsubscribed ` + user.Email + ` from our mailing list` + "\n"
	success += `If you'd like to resubscribe, please log in and go to "Email Settings" under "My Profile"`

	return success, nil
}

// Handler for /game path
func UnsubscribeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var msg string
		var err *ApplicationError

		switch r.Method {
		case "GET":
			msg, err = getUnsubscribe(r)

		default:
			msg = ""
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeNotFoundMethod)
		}
		WriteStringToPayload(w, r, msg, err)
	}
}
