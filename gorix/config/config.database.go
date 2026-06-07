package config

import (
	"fmt"
	"time"

	"github.com/Gromosome/gorix/gorix/core/database"
)

type DatabaseConfig struct {
	Driver                string
	DSN                   string
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime string
	ConnectionMaxIdleTime string
}

func (c Config) DatabaseConfigs() ([]database.Config, error) {
	configs := make(
		[]database.Config,
		0,
		len(c.Gorix.Databases),
	)

	for name, source := range c.Gorix.Databases {
		maxLifetime, err := parseDuration(
			source.ConnectionMaxLifetime,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"gorix config: invalid database %q connection-max-lifetime: %w",
				name,
				err,
			)
		}

		maxIdleTime, err := parseDuration(
			source.ConnectionMaxIdleTime,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"gorix config: invalid database %q connection-max-idle-time: %w",
				name,
				err,
			)
		}

		configs = append(
			configs,
			database.Config{
				Name:                  name,
				Driver:                source.Driver,
				DSN:                   source.DSN,
				MaxOpenConnections:    source.MaxOpenConnections,
				MaxIdleConnections:    source.MaxIdleConnections,
				ConnectionMaxLifetime: maxLifetime,
				ConnectionMaxIdleTime: maxIdleTime,
			},
		)
	}

	return configs, nil
}

func parseDuration(value string) (
	time.Duration,
	error,
) {
	if value == "" {
		return 0, nil
	}

	return time.ParseDuration(value)
}
