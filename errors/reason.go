package errors

type Reason string

const (
	ReasonInternal        Reason = "internal_server_error"
	ReasonNotFound        Reason = "not_found"
	ReasonBadRequest      Reason = "bad_request"
	ReasonUnauthorized    Reason = "unauthorized"
	ReasonForbidden       Reason = "forbidden"
	ReasonTooManyRequests Reason = "too_many_requests"
	ReasonConflict        Reason = "conflict"
	ReasonGatewayTimeout  Reason = "gateway_timeout"
	ReasonUnavailable     Reason = "unavailable"
)
