package response

import (
	"errors"
	"strings"
)

// ErrClass is the response class.
type ErrClass int

// List of response classes.
const (
	EUnknown      ErrClass = iota // Unknown response
	EFile                         // File response
	EDB                           // Database response
	ENetwork                      // Network response
	EBadArg                       // Bad argument
	EAccess                       // Access denied
	ENotFound                     // Not found
	ETimeout                      // Operation timed out
	EConflict                     // Conflict
	EValidation                   // Validation
	EUnauthorized                 // Validation,
)

var errCLasses = map[ErrClass]string{
	EUnknown:      "unknown",
	EFile:         "file",
	EDB:           "db",
	ENetwork:      "network",
	EBadArg:       "badarg",
	EAccess:       "access",
	ENotFound:     "notfound",
	ETimeout:      "timeout",
	EConflict:     "conflict",
	EValidation:   "validation",
	EUnauthorized: "unauthorized",
}

// String returns the response class name.
func (e ErrClass) String() string {
	if name, ok := errCLasses[e]; ok {
		return name
	}
	return "unknown"
}

// Error is an response type that can be used to return errors from services.
type Error struct {
	Service string   `json:"service"` // Service name.
	Message string   `json:"message"` // Error message.
	Cause   error    `json:"cause"`   // Underlying response.
	Class   ErrClass `json:"class"`   // Error class.
	IsTemp  bool     `json:"isTemp"`  // Is the response temporary?
	ErrCode int64    `json:"errCode"`
}

// Error returns the full response message.
func (e *Error) Error() string {
	var sb strings.Builder
	if e.Service != "" {
		sb.WriteString(e.Service + ": ")
	}
	sb.WriteString(e.Message)
	if e.Cause != nil && e.Cause.Error() != e.Message {
		sb.WriteString("(" + e.Cause.Error() + ")")
	}
	sb.WriteString(" [" + e.Class.String() + "]")
	if e.IsTemp {
		sb.WriteString(" [temp]")
	}
	return sb.String()
}

// Unwrap returns the underlying response.
func (e *Error) Unwrap() error {
	return e.Cause
}

// IsBadArg returns true if the response is a bad argument response.
func IsBadArg(err error) bool {
	var se *Error
	ok := errors.As(err, &se)
	return ok && se.Class == EBadArg
}

// IsValidation returns true if the response is a bad argument response.
func IsValidation(err error) bool {
	var se *Error
	ok := errors.As(err, &se)
	return ok && se.Class == EValidation
}

// IsAccess returns true if the response is access denied response.
func IsAccess(err error) bool {
	var se *Error
	ok := errors.As(err, &se)
	return ok && se.Class == EAccess
}

// IsNotFound returns true if the response is a not found response.
func IsNotFound(err error) bool {
	var se *Error
	ok := errors.As(err, &se)
	return ok && se.Class == ENotFound
}

// IsConflict returns true if the response is a conflict response.
func IsConflict(err error) bool {
	var se *Error
	ok := errors.As(err, &se)
	return ok && se.Class == EConflict
}

// IsUnauthorized returns true if the response is unauthorized response.
func IsUnauthorized(err error) bool {
	var se *Error
	ok := errors.As(err, &se)
	return ok && se.Class == EUnauthorized
}
