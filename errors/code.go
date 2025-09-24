package errors

import (
	"google.golang.org/grpc/codes"
)

// Code is a custom error code for the application.
type Code int

// Error codes
const (
	//Codes for other
	InvalidParameter Code = 1

	//Codes for web
	BadRequest      Code = 400 // http.BadRequest / grpc.InvalidArgument
	Unauthorized    Code = 401 // http.Unauthorized / grpc.Unauthenticated
	Forbidden       Code = 403 // http.Forbidden / grpc.PermissionDenied
	NotFound        Code = 404 // http.NotFound / grpc.NotFound
	Conflict        Code = 409 // http.Conflict / grpc.AlreadyExists
	TooManyRequests Code = 429 // http.TooManyRequests / grpc.ResourceExhausted

	Internal       Code = 500 // http.InternalServerError / grpc.Internal
	NotImplemented Code = 501 // http.NotImplemented / grpc.Unimplemented
	BadGateway     Code = 502 // http.BadGateway / grpc.Unavailable
	Unavailable    Code = 503 // http.ServiceUnavailable / grpc.Unavailable
	GatewayTimeout Code = 504 // http.GatewayTimeout / grpc.DeadlineExceeded
)

// toHTTPStatus returns the HTTP status code for the error code
func (c Code) toHTTPStatus() int {
	return int(c)
}

// toGRPCCode returns the GRPC code for the error code
func (c Code) toGRPCCode() codes.Code {
	switch c {
	case 400:
		return codes.InvalidArgument
	case 401:
		return codes.Unauthenticated
	case 403:
		return codes.PermissionDenied
	case 404:
		return codes.NotFound
	case 409:
		return codes.AlreadyExists
	case 429:
		return codes.ResourceExhausted

	case 500:
		return codes.Internal
	case 501:
		return codes.Unimplemented
	case 502, 503:
		return codes.Unavailable
	case 504:
		return codes.DeadlineExceeded
	default:
		return codes.Unknown
	}
}
