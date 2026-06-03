package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Gromosome/gorix/gorix/config"
)

type Config struct {
	Gorix GorixConfig
}

type GorixConfig struct {
	App AppConfig
}

type AppConfig struct {
	Prod *bool
	Host string
	Port int
}

func LoadConfig(root string) Config {
	cfg := DefaultConfig()

	path := filepath.Join(root, "application.yaml")

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}

	parsed := config.ParseYAML(string(data))

	cfg.Gorix.App.Prod = config.GetBoolPtr(parsed, "gorix.app.prod")
	cfg.Gorix.App.Host = config.GetString(parsed, "gorix.app.host", "0.0.0.0")
	cfg.Gorix.App.Port = config.GetInt(parsed, "gorix.app.port", 8080)

	return cfg
}

func DefaultConfig() Config {
	return Config{
		Gorix: GorixConfig{
			App: AppConfig{
				Prod: nil,
				Host: "0.0.0.0",
				Port: 8080,
			},
		},
	}
}

func (c Config) IsProd() bool {
	if c.Gorix.App.Prod == nil {
		return true
	}

	return *c.Gorix.App.Prod
}

func (c Config) Host() string {
	if c.Gorix.App.Host == "" {
		return "0.0.0.0"
	}

	return c.Gorix.App.Host
}

func (c Config) Port() int {
	if c.Gorix.App.Port == 0 {
		return 8080
	}

	return c.Gorix.App.Port
}

func (c Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host(), c.Port())
}
