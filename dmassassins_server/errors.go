package main

import (
	"database/sql"
	"fmt"
	"github.com/getsentry/raven-go"
	"log"
)

// Error codes are multiples of http codes for easy mapping
const (
	// 400 - Bad Request
	ErrCodeInvalidParameter = 40001
	ErrCodeMissingParameter = 40002
	ErrCodeInvalidHeader    = 40003
	ErrCodeMissingHeader    = 40004

	// 401 - Unauthorized
	ErrCodeNoSession           = 40100
	ErrCodeInvalidSecret       = 40101
	ErrCodeInvalidGamePassword = 40102
	ErrCodeInvalidFBToken      = 40120

	// 403 - Forbidden
	ErrCodePermissionDenied = 40300

	// 404 - Not Found -- Invalid Input
	ErrCodeInvalidMethod   = 40400
	ErrCodeInvalidEmail    = 40401
	ErrCodeInvalidUserId   = 40402
	ErrCodeInvalidUsername = 40403
	ErrCodeInvalidTeamId   = 40404
	ErrCodeInvalidGameId   = 40405

	// 404 - Not Found -- Valid Input
	ErrCodeNoGameMappings = 40420

	// 500 - Internal Server Error
	ErrCodeDatabase               = 50001
	ErrCodeDatabaseNoRowsAffected = 50002
	ErrCodeSession                = 50010 // Malformed Session
	ErrCodeWtf                    = 50069
)

// ApplicationError contains information about errors that arise while accessing resources.
type ApplicationError struct {
	Msg       string
	Err       error
	Code      int
	Exception *raven.Exception
}

// Error returns a human-readable representation of a ApplicationError.
func (err *ApplicationError) Error() (msg string) {
	return err.Msg
}

var sentryDSN string

// Creates a raven stacktrace
func trace() (stacktrace *raven.Stacktrace) {
	return raven.NewStacktrace(0, 2, nil)
}

// Converts an error to an Application Error with a user facing message
func NewApplicationError(msg string, err error, code int) (appErr *ApplicationError) {
	exception := raven.NewException(err, trace())
	return &ApplicationError{msg, err, code, exception}
}

// Determines if rows were affected in a sql result, reduces boilerplate on errors
func WereRowsAffected(res sql.Result) (appErr *ApplicationError) {
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabaseNoRowsAffected)
	}

	if rowsAffected == 0 {
		return NewApplicationError("Internal Error", err, ErrCodeDatabaseNoRowsAffected)
	}
	return nil
}

// LogWithSentry sends error report to sentry and records event id and error name to the logs
func LogWithSentry(appErr *ApplicationError, tags map[string]string, level raven.Severity, interfaces ...raven.Interface) {
	client, _ := raven.NewClient(Config.SentryDSN, tags)
	passthrough := append(interfaces, appErr.Exception)

	packet := raven.NewPacket(appErr.Error(), passthrough...)
	packet.Level = level
	packet.AddTags(tags)
	eventID, err := client.Capture(packet, nil)
	if err == nil {
		log.Print("Sentry failed to capture error below with message: ")
		log.Println(err)
	}
	message := fmt.Sprintf("Error event with id \"%s\" - %s", eventID, appErr.Error())
	log.Println(message)
}
