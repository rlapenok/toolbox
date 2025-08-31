package microservice

import "context"

// Gracefull - интерфейс для запуска и остановки сервиса
type Gracefull interface {
	// Name - имя сервиса
	Name() string
	// Address - адрес сервиса
	Address() string
	// Start - запуск сервиса
	Start() error
	// Stop - остановка сервиса
	Stop(ctx context.Context) error
}
