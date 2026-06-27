package document_driver_manager

import (
	"context"
	"fmt"
)

type Manager struct {
	client   Client
	adapter  Adapter
	database string
}

func NewManager(
	client Client,
	adapter Adapter,
	database string,
) *Manager {
	return &Manager{
		client:   client,
		adapter:  adapter,
		database: database,
	}
}

func Open(
	ctx context.Context,
	config Config,
) (*Manager, error) {
	adapter, err := Lookup(config.Driver)
	if err != nil {
		return nil, err
	}

	openCtx := ctx
	cancel := func() {}

	if config.PingTimeout > 0 {
		openCtx, cancel = context.WithTimeout(
			ctx,
			config.PingTimeout,
		)
	}
	defer cancel()

	client, err := adapter.Open(openCtx, config)
	if err != nil {
		return nil, adapter.Normalize(err)
	}

	if err := client.Ping(openCtx); err != nil {
		_ = client.Close(ctx)

		return nil, adapter.Normalize(err)
	}

	return &Manager{
		client:   client,
		adapter:  adapter,
		database: config.Database,
	}, nil
}

func (m *Manager) Driver() string {
	if m == nil || m.adapter == nil {
		return ""
	}

	return m.adapter.Name()
}

func (m *Manager) DatabaseName() string {
	if m == nil {
		return ""
	}

	return m.database
}

func (m *Manager) Normalize(err error) error {
	if err == nil {
		return nil
	}

	if m == nil || m.adapter == nil {
		return err
	}

	return m.adapter.Normalize(err)
}

func (m *Manager) Ping(ctx context.Context) error {
	if m == nil || m.client == nil {
		return fmt.Errorf(
			"gorix document-driver-manager: client is unavailable",
		)
	}

	return m.Normalize(
		m.client.Ping(ctx),
	)
}

func (m *Manager) Database() Database {
	if m == nil || m.client == nil {
		return nil
	}

	return m.client.Database(m.database)
}

func (m *Manager) Collection(name string) Collection {
	database := m.Database()
	if database == nil {
		return nil
	}

	return database.Collection(name)
}

func (m *Manager) Close() error {
	if m == nil || m.client == nil {
		return nil
	}

	return m.Normalize(
		m.client.Close(context.Background()),
	)
}
