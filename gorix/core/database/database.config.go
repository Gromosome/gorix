package database

import (
	"fmt"
	"time"
)

const DefaultConnectionName = "default"

type Config struct {
	Name                  string
	Driver                string
	DSN                   string
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
	ConnectionMaxIdleTime time.Duration
}

func (c Config) Normalize() Config {
	if c.Name == "" {
		c.Name = DefaultConnectionName
	}

	if c.MaxOpenConnections <= 0 {
		c.MaxOpenConnections = 25
	}

	if c.MaxIdleConnections < 0 {
		c.MaxIdleConnections = 0
	}

	if c.MaxIdleConnections == 0 {
		c.MaxIdleConnections = 10
	}

	if c.MaxIdleConnections > c.MaxOpenConnections {
		c.MaxIdleConnections = c.MaxOpenConnections
	}

	return c
}

func (c Config) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("gorix database: connection name is required")
	}

	if c.Driver == "" {
		return fmt.Errorf(
			"gorix database: driver is required for connection %q",
			c.Name,
		)
	}

	if c.DSN == "" {
		return fmt.Errorf(
			"gorix database: DSN is required for connection %q",
			c.Name,
		)
	}

	return nil
}
