package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type SupportPost struct {
	Email   string `json:"email"`
	Message string `json:"message"`
	Name    string `json:"name"`
	Subject string `json:"subject"`
}

// POST - Creates a support ticket
func postSupport(r *http.Request) (issueMap map[string]int, appErr *ApplicationError) {

	var supportPost SupportPost
	// Decode json
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&supportPost)
	if err != nil {
		msg := "Invalid JSON"
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeInvalidJSON)
	}

	// Check for missing parametsrs
	if supportPost.Name == "" {
		msg := "Missing Parameter: name"
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}
	if supportPost.Email == "" {
		msg := "Missing Parameter: email"
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}

	if supportPost.Message == "" {
		msg := "Missing Parameter: message"
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}

	if supportPost.Subject == "" {
		msg := "Missing Parameter: subject"
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}

	issueNum, appErr := postGithubIssue(supportPost.Subject, supportPost.Message, supportPost.Email, supportPost.Name)
	if appErr != nil {
		return nil, appErr
	}

	issueMap = map[string]int{
		"issue_num": issueNum,
	}

	return issueMap, nil

}

// Handler for /team path
func SupportHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var obj interface{}
		var err *ApplicationError

		switch r.Method {
		case "POST":
			obj, err = postSupport(r)
		default:
			obj = nil
			msg := "Not Found"
			tempErr := errors.New("Invalid Http Method")
			err = NewApplicationError(msg, tempErr, ErrCodeNotFoundMethod)
		}
		WriteObjToPayload(w, r, obj, err)
	}
}
