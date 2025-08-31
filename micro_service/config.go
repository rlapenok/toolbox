package microservice

// Config - интерфейс для конфигурации микросервиса
type Config interface {
	GetName() string
	EnableDefaultGinServer() bool
}
