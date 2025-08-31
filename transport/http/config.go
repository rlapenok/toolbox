package http

// Config - конфигурация HTTP сервера
type Config interface {
	// GetName - возвращает имя сервиса
	GetName() string
	// GetEnvironment - возвращает окружение в котором запущен сервис
	GetEnvironment() string
	// GetHost - возвращает хост сервиса
	GetHost() string
	// GetPort - возвращает порт сервиса
	GetPort() int
}
