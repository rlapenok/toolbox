package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware - middleware for logging requests
func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get the data from the context
		requestID := c.GetString("request_id")
		method := c.Request.Method
		path := c.Request.URL.Path
		query := c.Request.URL.Query().Encode()
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		requestSize := c.Request.ContentLength

		// add the data to the logger
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

		// log the start of the request
		midLogger.Info("request received")

		// get the start timestamp
		startTimestamp := time.Now()

		c.Next()

		// get the latency of the request, status and response size
		latency := time.Since(startTimestamp)
		status := c.Writer.Status()
		responseSize := c.Writer.Size()

		// add the latency, status and response size to the logger
		midLogger = midLogger.With(
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.Int("response_size", responseSize),
		)

		// log with appropriate level based on status code
		var logFunc func(string, ...zap.Field)
		message := "request completed"

		switch {
		case status >= 500:
			logFunc = midLogger.Error
		case status >= 400:
			logFunc = midLogger.Warn
		default:
			logFunc = midLogger.Info
		}

		// add error to log if present
		if len(c.Errors) > 0 && c.Errors.Last() != nil {
			logFunc(message, zap.Error(c.Errors.Last().Err))
		} else {
			logFunc(message)
		}
	}
}
