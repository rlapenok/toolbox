package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// EngineWrapper - функция для обертки движка Gin(добавление своих путей, middleware и т.д.)
type EngineWrapper func(engine *gin.Engine)

// GinServer - сервер HTTP на основе Gin
type GinServer struct {
	name    *string
	address *string
	port    int
	server  *http.Server
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

// Name - возвращает имя сервиса
func (s *GinServer) Name() string {
	return *s.name
}

// Address - возвращает адрес сервиса
func (s *GinServer) Address() string {
	return *s.address
}

// Start - запуск сервера
func (s *GinServer) Start() error {
	return s.server.ListenAndServe()
}

// Stop - остановка сервера
func (s *GinServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
