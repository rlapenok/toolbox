package postgres

import "time"

// ConnectionConfig - interface for connection configuration
type ConnectionConfig interface {
	GetHost() string
	GetPort() uint16
	GetUser() string
	GetPassword() string
	GetDatabase() string
	GetSchema() string
	GetSSLMode() string
	GetSSLCert() string
	GetSSLKey() string
	GetSSLRoot() string
}

// PoolConfig - interface for pool configuration
type PoolConfig interface {
	ConnectionConfig
	GetMinConns() int32
	GetMaxConns() int32
	GetMaxConnLifetime() time.Duration
	GetMaxConnIdleTime() time.Duration
	GetMaxConnKeepAliveTime() time.Duration
	GetMaxConnKeepAliveCount() int
	GetMaxConnKeepAliveInterval() time.Duration
}
