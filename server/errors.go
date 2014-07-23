package main

// Error codes are multiples of http codes for easy mapping
const (
	ErrCodeInvalidParameter = 40001
	ErrCodeMissingParameter = 40002

	ErrCodeNoSession     = 40100
	ErrCodeInvalidSecret = 40101

	ErrCodeInvalidMethod = 40400
	ErrCodeInvalidEmail  = 40401
	ErrCodeInvalidUserId = 40402

	ErrCodeDatabase = 50001
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

func NewApplicationError(msg string, err error, code int) *ApplicationError {

	return &ApplicationError{msg, err, code}
}

func CheckError(msg string, err error, code int) *ApplicationError {
	if err != nil {
		return &ApplicationError{msg, err, code}
	}
	return nil
}
