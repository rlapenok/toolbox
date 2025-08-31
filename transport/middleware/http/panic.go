package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rlapenok/toolbox/errors"
	"go.uber.org/zap"
)

// PanicMiddleware - middleware для обработки паники
func PanicMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Получаем данные из контекста
		requestID := c.GetString("request_id")
		method := c.Request.Method
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		requestSize := c.Request.ContentLength

		// Добавляем данные в логгер
		midLogger := logger.With(
			zap.String("request_id", requestID),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", ip),
			zap.String("user-agent", userAgent),
			zap.String("referer", referer),
			zap.Int64("request_size", requestSize),
		)

		startTimestamp := time.Now()

		defer func() {
			if r := recover(); r != nil {
				responseSize := c.Writer.Size()
				latency := time.Since(startTimestamp)
				details := errors.NewDetails()
				details.WithLocaleMessage("en-EN", "internal server error")

				err := errors.New(errors.Internal, "internal server error").
					WithDetails(details)

				httpStatus := err.ToHTTPStatus()
				code := err.Code()
				message := err.Message()
				reason := err.Reason()

				midLogger.With(
					zap.Int("status", httpStatus),
					zap.Duration("latency", latency),
					zap.Int("response_size", responseSize),
				).Error("request completed", zap.Any("panic", r))
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
