package database

import (
	"context"
	"database/sql"
	"fmt"
)

type Connection struct {
	name   string
	driver string
	db     *sql.DB
}

func Open(ctx context.Context, config Config) (*Connection, error) {
	config = config.Normalize()

	if err := config.Validate(); err != nil {
		return nil, err
	}

	if err := validateDriver(config.Driver); err != nil {
		return nil, err
	}

	db, err := sql.Open(config.Driver, config.DSN)
	if err != nil {
		return nil, fmt.Errorf(
			"gorix database: failed to open connection %q using driver %q: %w",
			config.Name,
			config.Driver,
			err,
		)
	}

	db.SetMaxOpenConns(config.MaxOpenConnections)
	db.SetMaxIdleConns(config.MaxIdleConnections)

	if config.ConnectionMaxLifetime > 0 {
		db.SetConnMaxLifetime(config.ConnectionMaxLifetime)
	}

	if config.ConnectionMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(config.ConnectionMaxIdleTime)
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()

		return nil, fmt.Errorf(
			"gorix database: failed to connect to %q: %w",
			config.Name,
			err,
		)
	}

	return &Connection{
		name:   config.Name,
		driver: config.Driver,
		db:     db,
	}, nil
}

func validateDriver(driver string) error {
	for _, registered := range sql.Drivers() {
		if registered == driver {
			return nil
		}
	}

	return fmt.Errorf(
		"gorix database: driver %q is not registered; add a blank import for the database/sql driver",
		driver,
	)
}

func (c *Connection) Name() string {
	return c.name
}

func (c *Connection) Driver() string {
	return c.driver
}

func (c *Connection) DB() *sql.DB {
	return c.db
}

func (c *Connection) Ping(ctx context.Context) error {
	if c == nil || c.db == nil {
		return fmt.Errorf("gorix database: connection is not initialized")
	}

	return c.db.PingContext(ctx)
}

func (c *Connection) Stats() sql.DBStats {
	if c == nil || c.db == nil {
		return sql.DBStats{}
	}

	return c.db.Stats()
}

func (c *Connection) Close() error {
	if c == nil || c.db == nil {
		return nil
	}

	return c.db.Close()
}
