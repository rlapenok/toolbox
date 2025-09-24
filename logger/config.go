package logger

// Config - config for logger
type Config interface {
	GetFormat() LogFormat
	GetLevel() string
	GetMeta() map[string]any
}
