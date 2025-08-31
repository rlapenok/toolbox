package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rlapenok/toolbox/errors"
)

// ErrorHandlerMiddleware - middleware for handling errors
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// call the next middleware
		c.Next()

		// get the last error
		ginError := c.Errors.Last()

		// if there is no error, return
		if ginError == nil {
			return
		}

		// try to cast the error to a toolbox error
		toolboxError, ok := ginError.Err.(*errors.Error)
		if !ok {
			// if the error is not a toolbox error, create a new toolbox error
			message := ginError.Err.Error()
			details := errors.NewDetails()
			details.WithLocaleMessage("en-EN", message)

			toolboxError = errors.New(errors.Internal, message).WithDetails(details)
		}

		httpStatus := toolboxError.ToHTTPStatus()
		code := toolboxError.Code()
		message := toolboxError.Message()
		details := toolboxError.Details()
		reason := toolboxError.Reason()

		c.AbortWithStatusJSON(httpStatus, gin.H{
			"code":    code,
			"message": message,
			"reason":  reason,
			"details": details,
		})
	}
}
