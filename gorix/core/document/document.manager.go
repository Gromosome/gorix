package document

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	docdriver "github.com/Gromosome/gorix/document-driver-manager"
	"github.com/Gromosome/gorix/gorix/config"
	gorixcontext "github.com/Gromosome/gorix/gorix/core/context"
)

type Connection struct {
	name     string
	config   config.DocumentConfig
	adapter  docdriver.Adapter
	client   docdriver.Client
	database docdriver.Database
}

func (c *Connection) Name() string {
	if c == nil {
		return ""
	}
	return c.name
}

func (c *Connection) Config() config.DocumentConfig {
	if c == nil {
		return config.DocumentConfig{}
	}
	return c.config
}

func (c *Connection) Client() docdriver.Client {
	if c == nil {
		return nil
	}
	return c.client
}

func (c *Connection) Database() docdriver.Database {
	if c == nil {
		return nil
	}
	return c.database
}

func (c *Connection) Collection(name string) docdriver.Collection {
	if c == nil || c.database == nil {
		return nil
	}
	return c.database.Collection(name)
}

func (c *Connection) Close(ctx context.Context) error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Close(ctx)
}

type Manager struct {
	mutex       sync.RWMutex
	connections map[string]*Connection
}

func NewManager() *Manager {
	return &Manager{
		connections: make(map[string]*Connection),
	}
}

func (m *Manager) Connect(
	ctx *gorixcontext.Context,
	config config.DocumentConfig,
) error {
	if m == nil {
		return fmt.Errorf("gorix document: manager cannot be nil")
	}

	config = config.Normalize()

	if config.Driver == "" {
		return fmt.Errorf(
			"gorix document: driver is required for connection %q",
			config.Name,
		)
	}

	if config.DSN == "" {
		return fmt.Errorf(
			"gorix document: DSN is required for connection %q",
			config.Name,
		)
	}
	if config.Database == "" {
		return fmt.Errorf(
			"gorix document: database is required for connection %q",
			config.Name,
		)
	}

	adapter, err := docdriver.Lookup(config.Driver)
	if err != nil {
		return err
	}

	var baseCtx context.Context = context.Background()
	if ctx != nil {
		baseCtx = ctx
	}

	openCtx := baseCtx
	cancel := func() {}

	pingTimeOut, _ := parseDuration(config.PingTimeout)
	if pingTimeOut > 0 {
		openCtx, cancel = context.WithTimeout(baseCtx, pingTimeOut)
	}
	defer cancel()

	client, err := adapter.Open(openCtx, config.DriverConfig())
	if err != nil {
		return adapter.Normalize(err)
	}

	if err := client.Ping(openCtx); err != nil {
		_ = client.Close(context.Background())
		return err
	}

	connection := &Connection{
		name:     config.Name,
		config:   config,
		adapter:  adapter,
		client:   client,
		database: client.Database(config.Database),
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.connections[config.Name]; exists {
		_ = client.Close(context.Background())

		return fmt.Errorf(
			"gorix document: connection %q is already registered",
			config.Name,
		)
	}

	m.connections[config.Name] = connection
	return nil
}

func (m *Manager) Connection(names ...string) (*Connection, error) {
	if m == nil {
		return nil, fmt.Errorf("gorix document: manager cannot be nil")
	}

	name := config.DefaultConnectionName
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	connection, exists := m.connections[name]
	if !exists {
		return nil, fmt.Errorf(
			"gorix document: connection %q is not registered",
			name,
		)
	}

	return connection, nil
}

func (m *Manager) Database(names ...string) (docdriver.Database, error) {
	connection, err := m.Connection(names...)
	if err != nil {
		return nil, err
	}

	return connection.Database(), nil
}

func (m *Manager) MustDatabase(names ...string) docdriver.Database {
	database, err := m.Database(names...)
	if err != nil {
		panic(err)
	}
	return database
}

func (m *Manager) Has(name string) bool {
	if m == nil {
		return false
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	_, exists := m.connections[name]
	return exists
}

func (m *Manager) Names() []string {
	if m == nil {
		return nil
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	names := make([]string, 0, len(m.connections))
	for name := range m.connections {
		names = append(names, name)
	}

	return names
}

func (m *Manager) Close() error {
	if m == nil {
		return nil
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	var closeErrors []error

	for name, connection := range m.connections {
		if err := connection.Close(context.Background()); err != nil {
			closeErrors = append(
				closeErrors,
				fmt.Errorf(
					"gorix document: failed to close connection %q: %w",
					name,
					err,
				),
			)
		}
	}

	m.connections = make(map[string]*Connection)

	return errors.Join(closeErrors...)
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
