package microservice

import (
	"context"

	"go.uber.org/zap"
)

// Gracefull - interface for start and stop service
type Gracefull interface {
	// Name - name of the service
	Name() string
	// Address - address of the service
	Address() string
	// Start - start the service
	Start() error
	// Stop - stop the service
	Stop(ctx context.Context) error

	// Logger - logger of the service
	Logger() *zap.Logger

	// WithLogger - set logger
	WithLogger(logger *zap.Logger) Gracefull
}
