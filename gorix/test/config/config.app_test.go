package config

import (
	"os"
	"path/filepath"
	"testing"

	config2 "github.com/Gromosome/gorix/gorix/config"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config2.DefaultConfig()
	if !cfg.IsProd() {
		t.Fatal("default config should be prod")
	}
	if cfg.Host() != "0.0.0.0" {
		t.Fatalf("unexpected default host: %s", cfg.Host())
	}
	if cfg.Port() != 8080 {
		t.Fatalf("unexpected default port: %d", cfg.Port())
	}
	if cfg.Address() != "0.0.0.0:8080" {
		t.Fatalf("unexpected default address: %s", cfg.Address())
	}
}

func TestTryLoadConfigReadsApplicationYAMLAndEnvFile(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "dev.env"), []byte("GORIX_TEST_APP_PORT=9090\nGORIX_TEST_DB_DSN=postgres://localhost/app\n"), 0o600); err != nil {
		t.Fatalf("failed to write env file: %v", err)
	}
	if err := os.Unsetenv("GORIX_TEST_APP_PORT"); err != nil {
		t.Fatalf("failed to unset env: %v", err)
	}
	if err := os.Unsetenv("GORIX_TEST_DB_DSN"); err != nil {
		t.Fatalf("failed to unset env: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Unsetenv("GORIX_TEST_APP_PORT")
		_ = os.Unsetenv("GORIX_TEST_DB_DSN")
	})

	content := []byte(`
env: dev
gorix:
  app:
    prod: false
    host: 127.0.0.1
    port: ${GORIX_TEST_APP_PORT}
  databases:
    primary:
      driver: postgres
      dsn: ${GORIX_TEST_DB_DSN}
      max-open-connections: 15
      max-idle-connections: 7
      connection-max-lifetime: 1h
      connection-max-idle-time: 5m
`)
	if err := os.WriteFile(filepath.Join(dir, "application.yaml"), content, 0o600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := config2.TryLoadConfig(dir)
	if err != nil {
		t.Fatalf("TryLoadConfig returned error: %v", err)
	}

	if cfg.Env != "dev" {
		t.Fatalf("unexpected env: %s", cfg.Env)
	}
	if cfg.IsProd() {
		t.Fatal("configured prod=false was not applied")
	}
	if cfg.Address() != "127.0.0.1:9090" {
		t.Fatalf("unexpected address: %s", cfg.Address())
	}
	if cfg.Gorix.Databases["primary"].DSN != "postgres://localhost/app" {
		t.Fatalf("unexpected database config: %#v", cfg.Gorix.Databases["primary"])
	}
}

func TestTryLoadConfigRejectsInvalidEnvironmentName(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "application.yaml"), []byte("env: ../prod\n"), 0o600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	if _, err := config2.TryLoadConfig(dir); err == nil {
		t.Fatal("expected invalid environment error")
	}
}
