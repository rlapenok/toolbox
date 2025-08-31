package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
)

// Error is a custom error type for the application.
type Error struct {
	// mapping to HTTP code and GRPC code
	errCode Code

	// developer message
	message string

	// reason of the error
	reason Reason

	// any details to be returned to the client
	details any
}

// implement the error interface
func (e *Error) Error() string {
	return fmt.Sprintf("code: %d, message: %s, details: %v", e.errCode, e.message, e.details)
}

// New creates a new error with the given code and message
func New(code Code, message string) *Error {
	return &Error{
		errCode: code,
		message: message,
		details: nil,
	}
}

// add details to the error
func (e *Error) WithDetails(details any) *Error {
	e.details = details

	return e
}

// add reason to the error
func (e *Error) WithReason(reason Reason) *Error {
	e.reason = reason

	return e
}

// return the HTTP status code
func (e *Error) ToHTTPStatus() int {
	return e.errCode.toHTTPStatus()
}

// return the GRPC codeq
func (e *Error) ToGRPCCode() codes.Code {
	return e.errCode.toGRPCCode()
}

// return the error code
func (e *Error) Code() Code {
	return e.errCode
}

// return the error message
func (e *Error) Message() string {
	return e.message
}

// return the error reason
func (e *Error) Reason() Reason {
	return e.reason
}

// return the error details
func (e *Error) Details() any {
	return e.details
}
