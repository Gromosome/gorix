package database

import (
	"fmt"
	"time"

	gorixcontext "github.com/Gromosome/gorix/gorix/core/context"
	sqldriver "github.com/Gromosome/gorix/sql-driver-manager"
)

type Stats struct {
	MaxOpenConnections int

	OpenConnections  int
	InUseConnections int
	IdleConnections  int

	WaitCount    int64
	WaitDuration time.Duration

	MaxIdleClosed     int64
	MaxIdleTimeClosed int64
	MaxLifetimeClosed int64
}

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

	driverManager, err := sqldriver.Open(
		ctx,
		sqldriver.Config{
			Driver:          config.Driver,
			DSN:             config.DSN,
			MaxOpenConns:    config.MaxOpenConnections,
			MaxIdleConns:    config.MaxIdleConnections,
			ConnMaxLifetime: config.ConnectionMaxLifetime,
			ConnMaxIdleTime: config.ConnectionMaxIdleTime,
			PingTimeout:     config.PingTimeout,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"gorix database: failed to open connection %q: %w",
			config.Name,
			err,
		)
	}

	return &Connection{
		name:   config.Name,
		driver: driverManager.Driver(),
		db: &DB{
			native: driverManager,
		},
	}, nil
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

func (c *Connection) Ping(
	ctx *gorixcontext.Context,
) error {
	if c == nil || c.db == nil {
		return fmt.Errorf(
			"gorix database: connection is unavailable",
		)
	}

	return c.db.Ping(ctx)
}

func (c *Connection) Stats() Stats {
	if c == nil ||
		c.db == nil ||
		c.db.native == nil {
		return Stats{}
	}

	stats := c.db.native.Stats()

	return Stats{
		MaxOpenConnections: stats.MaxOpenConnections,
		OpenConnections:    stats.OpenConnections,
		InUseConnections:   stats.InUse,
		IdleConnections:    stats.Idle,
		WaitCount:          stats.WaitCount,
		WaitDuration:       stats.WaitDuration,
		MaxIdleClosed:      stats.MaxIdleClosed,
		MaxIdleTimeClosed:  stats.MaxIdleTimeClosed,
		MaxLifetimeClosed:  stats.MaxLifetimeClosed,
	}
}

func (c *Connection) Close() error {
	if c == nil ||
		c.db == nil ||
		c.db.native == nil {
		return nil
	}

	return c.db.native.Close()
}
