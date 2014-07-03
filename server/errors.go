package main
 
// UserError contains error information presented to the users of our API
type UserError struct {
	Msg  string `json:"message"`
	Code int    `json:"error"`
}
 
 //error codes
const (
	ERROR_INVALID_SECRET = 40301
	ERROR_INVALID_EMAIL = 40401
	ERROR_INVALID_USER_ID = 40402
	ERROR_DATABASE = 50001

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