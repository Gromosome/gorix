package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Gromosome/gorix/gorix/config/yaml"
)

type Config struct {
	Env   string
	Gorix GorixConfig
}

type GorixConfig struct {
	App       AppConfig
	Databases map[string]DatabaseConfig
}

type AppConfig struct {
	Prod *bool
	Host string
	Port int
}

func LoadConfig(root string) Config {
	cfg, err := TryLoadConfig(root)
	if err != nil {
		panic(err)
	}

	return cfg
}

func TryLoadConfig(
	root string,
) (Config, error) {
	cfg := DefaultConfig()

	path := filepath.Join(
		root,
		"application.yaml",
	)

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}

		return Config{}, fmt.Errorf(
			"gorix config: failed to read %q: %w",
			path,
			err,
		)
	}

	content := string(data)

	bootstrap := yaml.ParseYAML(content)

	envName := strings.TrimSpace(
		yaml.GetString(
			bootstrap,
			"env",
			"",
		),
	)

	envFile := ".env"
	envFileRequired := false

	if envName != "" {
		if !yaml.IsValidEnvironmentName(envName) {
			return Config{}, fmt.Errorf(
				"gorix config: invalid environment name %q",
				envName,
			)
		}

		envFile = envName + ".env"
		envFileRequired = true
	}

	if err := yaml.LoadEnvFile(
		filepath.Join(root, envFile),
		envFileRequired,
	); err != nil {
		return Config{}, err
	}

	parsed := yaml.ParseYAML(content)

	cfg = buildConfig(parsed)
	cfg.Env = envName

	return cfg, nil
}

func buildConfig(
	parsed map[string]yaml.YAMLValue,
) Config {
	cfg := DefaultConfig()

	cfg.Env = strings.TrimSpace(
		yaml.GetString(
			parsed,
			"env",
			"",
		),
	)

	cfg.Gorix.App = AppConfig{
		Prod: yaml.GetBoolPtr(
			parsed,
			"gorix.app.prod",
		),
		Host: yaml.GetString(
			parsed,
			"gorix.app.host",
			"0.0.0.0",
		),
		Port: yaml.GetInt(
			parsed,
			"gorix.app.port",
			8080,
		),
	}

	cfg.Gorix.Databases =
		loadDatabaseConfigs(parsed)

	return cfg
}

func DefaultConfig() Config {
	return Config{
		Env: "",
		Gorix: GorixConfig{
			App: AppConfig{
				Prod: nil,
				Host: "0.0.0.0",
				Port: 8080,
			},
			Databases: make(
				map[string]DatabaseConfig,
			),
		},
	}
}
func loadDatabaseConfigs(
	parsed map[string]yaml.YAMLValue,
) map[string]DatabaseConfig {
	databaseSources, ok := yaml.GetMap(
		parsed,
		"gorix.databases",
	)
	if !ok {
		return make(map[string]DatabaseConfig)
	}

	databases := make(
		map[string]DatabaseConfig,
		len(databaseSources),
	)

	for name, value := range databaseSources {
		source, ok := value.(map[string]yaml.YAMLValue)
		if !ok {
			continue
		}

		databases[name] = DatabaseConfig{
			Driver:                yaml.GetString(source, "driver", ""),
			DSN:                   yaml.GetString(source, "dsn", ""),
			MaxOpenConnections:    yaml.GetInt(source, "max-open-connections", 0),
			MaxIdleConnections:    yaml.GetInt(source, "max-idle-connections", 0),
			ConnectionMaxLifetime: yaml.GetString(source, "connection-max-lifetime", ""),
			ConnectionMaxIdleTime: yaml.GetString(source, "connection-max-idle-time", ""),
			PingTimeout:           yaml.GetString(source, "ping-timeout", ""),
		}
	}

	return databases
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
	return fmt.Sprintf(
		"%s:%d",
		c.Host(),
		c.Port(),
	)
}
