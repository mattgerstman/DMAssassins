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

	// 404 - Not Found -- Valid Input
	ErrCodeNoGameMappings = 40420

	// 500 - Internal Server Error
	ErrCodeDatabase               = 50001
	ErrCodeDatabaseNoRowsAffected = 50002
	ErrCodeSession                = 50010 // Malformed Session
)

// ApplicationError contains information about errors that arise while accessing resources.
type ApplicationError struct {
	Msg       string
	Err       error
	Code      int
	Exception *raven.Exception
}

// Error returns a human-readable representation of a ApplicationError.
func (err *ApplicationError) Error() string {
	return err.Msg
}

var sentryDSN string

func trace() *raven.Stacktrace {
	return raven.NewStacktrace(0, 2, nil)
}

func NewApplicationError(msg string, err error, code int) *ApplicationError {
	exception := raven.NewException(err, trace())
	return &ApplicationError{msg, err, code, exception}
}

func CheckError(msg string, err error, code int) *ApplicationError {
	if err != nil {
		return NewApplicationError(msg, err, code)
	}
	return nil
}

func WereRowsAffected(res sql.Result) *ApplicationError {
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabaseNoRowsAffected)
	}

	if rowsAffected == 0 {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
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
	fmt.Println(err)
	message := fmt.Sprintf("Error event with id \"%s\" - %s", eventID, appErr.Error())
	log.Println(message)
}
