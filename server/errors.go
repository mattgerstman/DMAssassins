package main

import (
	"errors"
)

// UserError contains error information presented to the users of our API
type UserError struct {
	Msg  string `json:"message"`
	Code int    `json:"error"`
}

//error codes
const (
	ERROR_INVALID_PARAMETER = 40001
	ERROR_MISSING_PARAMETER = 40002
	ERROR_NO_SESSION		= 40100
	ERROR_INVALID_SECRET    = 40101
	ERROR_INVALID_METHOD    = 40400
	ERROR_INVALID_EMAIL     = 40401
	ERROR_INVALID_USER_ID   = 40402
	ERROR_DATABASE          = 50001
)

// ApplicationError contains information about errors that arise while accessing resources.
type ApplicationError struct {
	Msg  string
	Err  error
	Code int
}

// Error returns a human-readable representation of a ApplicationError.
func (err *ApplicationError) Error() string {
	return err.Msg
}

// UserError return a user-facing error
func (rerr *ApplicationError) UserError() *UserError {
	return &UserError{Msg: rerr.Msg, Code: rerr.Code}
}

func NewSimpleApplicationError(msg string, code int) *ApplicationError {
	return &ApplicationError{msg, errors.New(msg), code}
}

func CheckError(msg string, err error, code int) *ApplicationError {
	if err != nil {
		return &ApplicationError{msg, err, code}
	}
	return nil
}
