package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// EngineWrapper - функция для обертки движка Gin(добавление своих путей, middleware и т.д.)
type EngineWrapper func(engine *gin.Engine)

// GinServer - сервер HTTP на основе Gin
type GinServer struct {
	name    *string
	address *string
	port    int
	server  *http.Server
	logger  *zap.Logger
}

// NewGinServer - создание нового сервера HTTP на основе Gin
func NewGinServer(config Config, wrapper EngineWrapper) *GinServer {
	if config.GetEnvironment() == "production" || config.GetEnvironment() == "prod" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	if wrapper != nil {
		wrapper(engine)
	}

	port := config.GetPort()
	host := config.GetHost()
	name := config.GetName()
	address := fmt.Sprintf("%s:%d", host, port)

	server := &http.Server{
		Addr:    address,
		Handler: engine,
	}

	return &GinServer{
		name:    &name,
		address: &address,
		port:    port,
		server:  server,
	}
}

//===============================================
// Gracefull
//===============================================

// Name - return name of the service
func (s *GinServer) Name() string {
	return *s.name
}

// Address - return address of the service
func (s *GinServer) Address() string {
	return *s.address
}

// Start - start the service
func (s *GinServer) Start() error {
	return s.server.ListenAndServe()
}

// Stop - stop the service
func (s *GinServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// WithLogger - set logger to the service
func (s *GinServer) WithLogger(logger *zap.Logger) *GinServer {
	s.logger = logger
	return s
}

// Logger - return logger of the service
func (s *GinServer) Logger() *zap.Logger {
	return s.logger
}
