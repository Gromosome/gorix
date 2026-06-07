package database

import (
	"errors"
	"fmt"
	"sync"

	gorixcontext "github.com/Gromosome/gorix/gorix/core/context"
)

type Manager struct {
	mutex       sync.RWMutex
	connections map[string]*Connection
}

func NewManager() *Manager {
	return &Manager{
		connections: make(map[string]*Connection),
	}
}

func (m *Manager) Connect(ctx *gorixcontext.Context, config Config) error {
	config = config.Normalize()

	connection, err := Open(ctx, config)
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.connections[config.Name]; exists {
		_ = connection.Close()

		return fmt.Errorf(
			"gorix database: connection %q is already registered",
			config.Name,
		)
	}

	m.connections[config.Name] = connection
	return nil
}

func (m *Manager) Connection(names ...string) (*Connection, error) {
	name := DefaultConnectionName

	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	connection, exists := m.connections[name]
	if !exists {
		return nil, fmt.Errorf(
			"gorix database: connection %q is not registered",
			name,
		)
	}

	return connection, nil
}

func (m *Manager) DB(
	names ...string,
) (*DB, error) {
	connection, err := m.Connection(names...)
	if err != nil {
		return nil, err
	}

	return connection.DB(), nil
}

func (m *Manager) MustDB(names ...string) *DB {
	db, err := m.DB(names...)
	if err != nil {
		panic(err)
	}

	return db
}

func (m *Manager) Has(name string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	_, exists := m.connections[name]
	return exists
}

func (m *Manager) Names() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	names := make([]string, 0, len(m.connections))

	for name := range m.connections {
		names = append(names, name)
	}

	return names
}

func (m *Manager) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var closeErrors []error

	for name, connection := range m.connections {
		if err := connection.Close(); err != nil {
			closeErrors = append(
				closeErrors,
				fmt.Errorf(
					"gorix database: failed to close connection %q: %w",
					name,
					err,
				),
			)
		}
	}

	m.connections = make(map[string]*Connection)

	return errors.Join(closeErrors...)
}
