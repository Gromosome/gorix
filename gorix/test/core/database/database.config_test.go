package database

import (
	"testing"

	database2 "github.com/Gromosome/gorix/gorix/core/database"
)

func TestConfigNormalizeAppliesDefaultsAndBounds(t *testing.T) {
	cfg := database2.Config{
		MaxOpenConnections: -1,
		MaxIdleConnections: 50,
	}.Normalize()

	if cfg.Name != database2.DefaultConnectionName {
		t.Fatalf("unexpected default name: %s", cfg.Name)
	}
	if cfg.MaxOpenConnections != 25 {
		t.Fatalf("unexpected max open connections: %d", cfg.MaxOpenConnections)
	}
	if cfg.MaxIdleConnections != 25 {
		t.Fatalf("max idle should be capped to max open, got %d", cfg.MaxIdleConnections)
	}
}

func TestConfigValidate(t *testing.T) {
	valid := database2.Config{Name: "primary", Driver: "postgres", DSN: "dsn"}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid config returned error: %v", err)
	}

	tests := []database2.Config{
		{},
		{Name: "primary"},
		{Name: "primary", Driver: "postgres"},
	}
	for _, cfg := range tests {
		if err := cfg.Validate(); err == nil {
			t.Fatalf("expected validation error for %#v", cfg)
		}
	}
}
