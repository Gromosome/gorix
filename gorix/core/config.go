package core

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Gorix GorixConfig `yaml:"gorix"`
}

type GorixConfig struct {
	App AppConfig `yaml:"app"`
}

type AppConfig struct {
	Prod *bool  `yaml:"prod"`
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func LoadConfig(root string) Config {
	config := Config{
		Gorix: GorixConfig{
			App: AppConfig{
				Prod: nil,
				Host: "0.0.0.0",
				Port: 8080,
			},
		},
	}

	path := filepath.Join(root, "application.yaml")

	data, err := os.ReadFile(path)
	if err != nil {
		return config
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return config
	}

	if config.Gorix.App.Host == "" {
		config.Gorix.App.Host = "0.0.0.0"
	}

	if config.Gorix.App.Port == 0 {
		config.Gorix.App.Port = 8080
	}

	return config
}

func (c Config) IsProd() bool {
	if c.Gorix.App.Prod == nil {
		return false
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
