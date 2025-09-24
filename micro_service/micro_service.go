package microservice

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rlapenok/toolbox/logger"
	"go.uber.org/zap"
)

// MicroService - composite structure for microservice
// contains HTTP server, GRPC server, database connection, logger, configuration, etc.
type MicroService struct {
	name       *string
	gracefulls []Gracefull
	logger     *zap.Logger
}

// NewMicroService - create new microservice from config
func NewMicroService(config Config) *MicroService {
	logger := defaultLogger()

	gracefulls := []Gracefull{}

	name := config.GetName()

	return &MicroService{
		name:       &name,
		gracefulls: gracefulls,
		logger:     logger,
	}
}

// WithLogger - return new microservice with new logger
func (s *MicroService) WithLogger(config logger.Config) *MicroService {
	logger, err := logger.New(config)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	if s.logger != nil {
		if err := s.logger.Sync(); err != nil {
			log.Fatalf("failed to sync logger: %v", err)
		}

	}

	return &MicroService{
		gracefulls: s.gracefulls,
		logger:     logger,
	}
}

// WithGracefull - return new microservice with new gracefull
func (s *MicroService) WithGracefull(gracefull Gracefull) *MicroService {
	return &MicroService{
		gracefulls: append(s.gracefulls, gracefull),
		logger:     s.logger,
	}
}

// GetLogger - return logger
func (s *MicroService) GetLogger() *zap.Logger {
	return s.logger
}

// Run - run microservice
func (s *MicroService) Run(ctx context.Context) {
	if len(s.gracefulls) == 0 {
		s.logger.Error("no gracefulls to start")

		return
	}

	s.logger.Info("starting microservice...")

	errChan := make(chan error, 1)

	for _, gracefull := range s.gracefulls {
		go func() {
			s.logger.Info("starting gracefull",
				zap.String("name", gracefull.Name()),
				zap.String("address", gracefull.Address()),
			)

			if err := gracefull.Start(); err != nil && err != http.ErrServerClosed {
				s.logger.Error("failed to start gracefull",
					zap.String("name", gracefull.Name()),
					zap.String("address", gracefull.Address()),
					zap.Error(err),
				)
				errChan <- err
			}
		}()
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	select {
	//TODO: add graceful shutdown
	case <-ctx.Done():
	case <-errChan:
		for _, gracefull := range s.gracefulls {
			if err := gracefull.Stop(ctx); err != nil {
				s.logger.Error("failed to stop gracefull",
					zap.String("name", gracefull.Name()),
					zap.String("address", gracefull.Address()),
					zap.Error(err),
				)
			}
		}
	case <-signals:
		s.logger.Info("received shutdown signal")
		for _, gracefull := range s.gracefulls {
			if err := gracefull.Stop(context.Background()); err != nil {
				s.logger.Error("failed to stop gracefull",
					zap.String("name", gracefull.Name()),
					zap.String("address", gracefull.Address()),
					zap.Error(err),
				)
			}

			s.logger.Info("gracefull stopped",
				zap.String("name", gracefull.Name()),
				zap.String("address", gracefull.Address()),
			)
		}
	}

	s.logger.Info("microservice stopped")

	s.logger.Sync()
}
