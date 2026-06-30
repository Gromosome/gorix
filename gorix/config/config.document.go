package config

import (
	"strings"
	"time"

	docdriver "github.com/Gromosome/gorix/document-driver-manager"
	"github.com/Gromosome/gorix/gorix/config/yaml"
)

const DefaultConnectionName = "default"

type DocumentConfig struct {
	Name        string `yaml:"name" mapstructure:"name"`
	Driver      string `yaml:"driver" mapstructure:"driver"`
	DSN         string `yaml:"dsn" mapstructure:"dsn"`
	Database    string `yaml:"database" mapstructure:"database"`
	PingTimeout string `yaml:"ping-timeout" mapstructure:"ping-timeout"`
}

func (c Config) DocumentConfigs() ([]DocumentConfig, error) {
	configs := make(
		[]DocumentConfig,
		0,
		len(c.Gorix.Documents),
	)

	for name, source := range c.Gorix.Documents {

		configs = append(
			configs,
			DocumentConfig{
				Name:        name,
				Driver:      source.Driver,
				DSN:         source.DSN,
				Database:    source.Database,
				PingTimeout: source.PingTimeout,
			},
		)
	}

	return configs, nil
}

func loadDocumentConfigs(
	parsed map[string]yaml.YAMLValue,
) map[string]DocumentConfig {
	documentSources, ok := yaml.GetMap(
		parsed,
		"gorix.documents",
	)
	if !ok {
		return make(map[string]DocumentConfig)
	}

	documents := make(
		map[string]DocumentConfig,
		len(documentSources),
	)

	for name, value := range documentSources {
		source, ok := value.(map[string]yaml.YAMLValue)
		if !ok {
			continue
		}

		documents[name] = DocumentConfig{
			Driver: yaml.GetString(
				source,
				"driver",
				"",
			),
			DSN: yaml.GetString(
				source,
				"dsn",
				"",
			),
			Database: yaml.GetString(
				source,
				"database",
				"",
			),
			PingTimeout: yaml.GetString(
				source,
				"ping-timeout",
				"",
			),
		}
	}

	return documents
}

func (c DocumentConfig) Normalize() DocumentConfig {
	c.Name = strings.TrimSpace(c.Name)
	c.Driver = strings.ToLower(strings.TrimSpace(c.Driver))
	c.DSN = strings.TrimSpace(c.DSN)
	c.Database = strings.TrimSpace(c.Database)

	if c.Name == "" {
		c.Name = DefaultConnectionName
	}

	if c.Database == "" {
		c.Database = c.Name
	}
	pingTimeOut, _ := parseDuration(c.PingTimeout)
	if pingTimeOut <= 0 {
		pingTimeOut = 5 * time.Second
	}

	return c
}

func (c DocumentConfig) DriverConfig() docdriver.Config {
	c = c.Normalize()
	pingTimeOut, _ := parseDuration(c.PingTimeout)
	return docdriver.Config{
		Driver:      c.Driver,
		DSN:         c.DSN,
		Database:    c.Database,
		PingTimeout: pingTimeOut,
	}
}
