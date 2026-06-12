package config

import (
	"testing"
	"time"

	config2 "github.com/Gromosome/gorix/gorix/config"
)

func TestDatabaseConfigsParsesDurations(t *testing.T) {
	cfg := config2.DefaultConfig()
	cfg.Gorix.Databases["primary"] = config2.DatabaseConfig{
		Driver:                "postgres",
		DSN:                   "postgres://localhost/app",
		MaxOpenConnections:    5,
		MaxIdleConnections:    2,
		ConnectionMaxLifetime: "1h",
		ConnectionMaxIdleTime: "5m",
		PingTimeout:           "2s",
	}

	configs, err := cfg.DatabaseConfigs()
	if err != nil {
		t.Fatalf("DatabaseConfigs returned error: %v", err)
	}
	if len(configs) != 1 {
		t.Fatalf("expected one config, got %d", len(configs))
	}
	if configs[0].ConnectionMaxLifetime != time.Hour {
		t.Fatalf("unexpected max lifetime: %v", configs[0].ConnectionMaxLifetime)
	}
	if configs[0].PingTimeout != 2*time.Second {
		t.Fatalf("unexpected ping timeout: %v", configs[0].PingTimeout)
	}
}

func TestDatabaseConfigsRejectsInvalidDuration(t *testing.T) {
	cfg := config2.DefaultConfig()
	cfg.Gorix.Databases["primary"] = config2.DatabaseConfig{
		ConnectionMaxLifetime: "invalid",
	}

	if _, err := cfg.DatabaseConfigs(); err == nil {
		t.Fatal("expected invalid duration error")
	}
}
