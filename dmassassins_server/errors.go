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
	ErrCodeBadRequestWTF    = 40000
	ErrCodeInvalidParameter = 40001
	ErrCodeMissingParameter = 40002
	ErrCodeInvalidHeader    = 40003
	ErrCodeMissingHeader    = 40004
	ErrCodeInvalidUUID      = 40005
	ErrCodeInvalidJSON      = 40006
	ErrCodeNeedMorePlayers  = 40007

	// 401 - Unauthorized
	ErrCodeUnauthorizedWTF     = 40100
	ErrCodeNoSession           = 40101
	ErrCodeInvalidSecret       = 40102
	ErrCodeInvalidGamePassword = 40103
	ErrCodeInvalidFBToken      = 40120

	// 402 - Payment Required
	// Who the hell uses that?

	// 403 - Forbidden
	ErrCodeForbiddenWTF     = 40300
	ErrCodePermissionDenied = 40301

	// 404 - Not Found -- Invalid Input
	ErrCodeNotFoundWTF         = 40400
	ErrCodeNotFoundMethod      = 40401
	ErrCodeNotFoundEmail       = 40402
	ErrCodeNotFoundUsername    = 40403
	ErrCodeNotFoundUserId      = 40404
	ErrCodeNotFoundTeamId      = 40405
	ErrCodeNotFoundGameId      = 40406
	ErrCodeNotFoundGameMapping = 40407
	ErrCodeNotFoundTarget      = 40408

	// 404 - Not Found -- Valid Input
	ErrCodeNoGameMappings = 40420
	ErrCodeNoUsers        = 40421
	ErrCodeNoTeams        = 40422

	// 500 - Internal Server Error
	ErrCodeInternalServerWTF      = 50000
	ErrCodeDatabase               = 50001
	ErrCodeDatabaseNoRowsAffected = 50002
	ErrCodeFile                   = 50003
	ErrCodeEmail                  = 50004
	ErrCodeBadTemplate            = 50005
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
func (err *ApplicationError) Error() (msg string) {
	return err.Msg
}

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
