package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rlapenok/toolbox/errors"
	"go.uber.org/zap"
)

// PanicMiddleware - middleware для обработки паники
func PanicMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		defer func() {
			if r := recover(); r != nil {
				details := errors.NewDetails()
				details.WithLocaleMessage("en-EN", "internal server error")

				err := errors.New(errors.Internal, "internal server error").
					WithDetails(details)

				httpStatus := err.ToHTTPStatus()
				code := err.Code()
				message := err.Message()
				reason := err.Reason()

				logger.Error("PANIC", zap.Any("panic", r))

				c.AbortWithStatusJSON(httpStatus, gin.H{
					"code":    code,
					"message": message,
					"reason":  reason,
					"details": details,
				})
			}
		}()

		c.Next()
	}
}
