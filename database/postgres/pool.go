package postgres

import (
	"context"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rlapenok/toolbox/database"

	pgxMigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Pool - connection pool to PostgreSQL
type Pool struct {
	pool *pgxpool.Pool
}

// NewPool - create new pool
func NewPool(ctx context.Context, config PoolConfig) (*Pool, error) {

	// Parse config
	poolConfig, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, MapError(err)
	}

	// Get connection config
	conn := &poolConfig.ConnConfig.Config

	// Set connection config
	conn.Host = config.GetHost()
	conn.Port = config.GetPort()
	conn.User = config.GetUser()
	conn.Password = config.GetPassword()
	conn.Database = config.GetDatabase()

	// Set search path
	if config.GetSchema() != "" {
		conn.RuntimeParams["search_path"] = config.GetSchema()
	}

	// Set TLS config
	if config.GetSSLMode() != "disable" && config.GetSSLCert() != "" && config.GetSSLKey() != "" {
		tlsCfg, err := database.BuildTLSConfig(config.GetSSLRoot(), config.GetSSLCert(), config.GetSSLKey())
		if err != nil {
			return nil, MapError(err)
		}
		conn.TLSConfig = tlsCfg
	}

	// Set pool config
	poolConfig.MinConns = config.GetMinConns()
	poolConfig.MaxConns = config.GetMaxConns()
	poolConfig.MaxConnLifetime = config.GetMaxConnLifetime()
	poolConfig.MaxConnIdleTime = config.GetMaxConnIdleTime()

	// Create pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, MapError(err)
	}

	return &Pool{pool: pool}, nil
}

func (p *Pool) Migrate(ctx context.Context, migrationsPath string) error {
	conn := stdlib.OpenDB(*p.pool.Config().ConnConfig)

	driver, err := pgxMigrate.WithInstance(conn, &pgxMigrate.Config{})
	if err != nil {
		return MapError(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"pgx",
		driver,
	)
	if err != nil {
		return MapError(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return MapError(err)
	}

	if err := conn.Close(); err != nil {
		return MapError(err)
	}

	return nil
}

// Pgx - get pgx pool
func (p *Pool) Pgx() *pgxpool.Pool {
	return p.pool
}

// Close - close pool
func (p *Pool) Close() {
	p.pool.Close()
}
