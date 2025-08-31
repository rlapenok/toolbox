package microservice

import (
	"context"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rlapenok/toolbox/errors"
	"github.com/rlapenok/toolbox/transport/http"
	middlewareHTTP "github.com/rlapenok/toolbox/transport/middleware/http"
	"go.uber.org/zap"
)

type httpConfig struct {
	environment string
	host        string
	port        int
}

func (c *httpConfig) GetEnvironment() string {
	return c.environment
}

func (c *httpConfig) GetHost() string {
	return c.host
}

func (c *httpConfig) GetPort() int {
	return c.port
}

func (c *httpConfig) GetName() string {
	return "default gin server"
}

func defaultGinServer(logger *zap.Logger) *http.GinServer {
	config := &httpConfig{
		environment: "development",
		host:        "0.0.0.0",
		port:        8080,
	}

	wrapper := func(engine *gin.Engine) {

		engine.Use(middlewareHTTP.RequestIDMiddleware())
		engine.Use(middlewareHTTP.PanicMiddleware(logger))
		engine.Use(middlewareHTTP.LoggerMiddleware(logger))
		engine.Use(middlewareHTTP.ErrorHandlerMiddleware())

		engine.GET("/panic", panicHandler)

		engine.GET("/livez", livezHandler)

		engine.GET("/readyz", readyzHandler)

	}

	return http.NewGinServer(config, wrapper)
}

func panicHandler(c *gin.Context) {
	panic("test panic")
}

func livezHandler(c *gin.Context) {
	c.AbortWithStatusJSON(200, gin.H{
		"status": "ok",
	})
}

func readyzHandler(c *gin.Context) {
	reqCtx := c.Request.Context()
	ctx, cancel := context.WithTimeout(reqCtx, 500*time.Millisecond)
	defer cancel()

	sleep := rand.Intn(1000)

	select {
	case <-ctx.Done():
		details := errors.NewDetails()
		details.WithLocaleMessage("en-EN", "not ready")

		err := errors.New(errors.Unavailable, "not ready").
			WithReason(errors.ReasonReadz).
			WithDetails(details)

		http.ReturnGinError(c, err, nil)
	case <-time.After(time.Duration(sleep) * time.Millisecond):
		c.AbortWithStatusJSON(200, gin.H{
			"status": "ready",
		})
	}
}
