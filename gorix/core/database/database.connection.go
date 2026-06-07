package database

import (
	"database/sql"
	"fmt"

	gorixcontext "github.com/Gromosome/gorix/gorix/core/context"
)

type Connection struct {
	name   string
	driver string
	db     *DB
}

func Open(
	ctx *gorixcontext.Context,
	config Config,
) (*Connection, error) {
	config = config.Normalize()

	if err := config.Validate(); err != nil {
		return nil, err
	}

	if ctx == nil {
		return nil, fmt.Errorf(
			"gorix database: context cannot be nil",
		)
	}

	if err := validateDriver(config.Driver); err != nil {
		return nil, err
	}

	nativeDB, err := sql.Open(
		config.Driver,
		config.DSN,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"gorix database: failed to open connection %q using driver %q: %w",
			config.Name,
			config.Driver,
			err,
		)
	}

	nativeDB.SetMaxOpenConns(
		config.MaxOpenConnections,
	)
	nativeDB.SetMaxIdleConns(
		config.MaxIdleConnections,
	)

	if config.ConnectionMaxLifetime > 0 {
		nativeDB.SetConnMaxLifetime(
			config.ConnectionMaxLifetime,
		)
	}

	if config.ConnectionMaxIdleTime > 0 {
		nativeDB.SetConnMaxIdleTime(
			config.ConnectionMaxIdleTime,
		)
	}

	db := &DB{
		native: nativeDB,
	}

	if err := db.Ping(ctx); err != nil {
		_ = nativeDB.Close()

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

func (c *Connection) DB() *DB {
	return c.db
}

func (c *Connection) Ping(ctx *gorixcontext.Context) error {
	return c.db.Ping(ctx)
}

func (c *Connection) Stats() sql.DBStats {
	if c == nil || c.db == nil {
		return sql.DBStats{}
	}

	return c.db.native.Stats()
}

func (c *Connection) Close() error {
	if c == nil ||
		c.db == nil ||
		c.db.native == nil {
		return nil
	}

	return c.db.native.Close()
}
