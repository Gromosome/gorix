package sql_driver_manager

import (
	"context"
	"database/sql"
	"errors"
)

type Manager struct {
	db      *sql.DB
	adapter Adapter
}

func Open(ctx context.Context, config Config) (*Manager, error) {
	adapter, err := Lookup(config.Driver)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(adapter.SQLDriverName(), config.DSN)
	if err != nil {
		return nil, adapter.Normalize(err)
	}

	if config.MaxOpenConns > 0 {
		db.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.MaxIdleConns > 0 {
		db.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(config.ConnMaxLifetime)
	}
	if config.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(config.ConnMaxIdleTime)
	}

	pingContext := ctx
	cancel := func() {}
	if config.PingTimeout > 0 {
		pingContext, cancel = context.WithTimeout(ctx, config.PingTimeout)
	}
	defer cancel()

	if err = db.PingContext(pingContext); err != nil {
		_ = db.Close()
		return nil, adapter.Normalize(err)
	}

	return &Manager{
		db:      db,
		adapter: adapter,
	}, nil
}

func (m *Manager) Driver() string {
	return m.adapter.Name()
}

func (m *Manager) Normalize(err error) error {
	if err == nil {
		return nil
	}
	return m.adapter.Normalize(err)
}

func (m *Manager) Close() error {
	return m.Normalize(m.db.Close())
}

func (m *Manager) Stats() sql.DBStats {
	return m.db.Stats()
}

func (m *Manager) PingContext(ctx context.Context) error {
	return m.Normalize(m.db.PingContext(ctx))
}

func (m *Manager) ExecContext(
	ctx context.Context,
	query string,
	args ...any,
) (sql.Result, error) {
	result, err := m.db.ExecContext(ctx, query, args...)
	return result, m.Normalize(err)
}

func (m *Manager) QueryContext(
	ctx context.Context,
	query string,
	args ...any,
) (*Rows, error) {
	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, m.Normalize(err)
	}

	return &Rows{
		Rows:      rows,
		normalize: m.Normalize,
	}, nil
}

func (m *Manager) QueryRowContext(
	ctx context.Context,
	query string,
	args ...any,
) *Row {
	return &Row{
		row:       m.db.QueryRowContext(ctx, query, args...),
		normalize: m.Normalize,
	}
}

func (m *Manager) PrepareContext(
	ctx context.Context,
	query string,
) (*Stmt, error) {
	statement, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, m.Normalize(err)
	}

	return &Stmt{
		stmt:      statement,
		normalize: m.Normalize,
	}, nil
}

func (m *Manager) BeginTx(
	ctx context.Context,
	options *sql.TxOptions,
) (*Tx, error) {
	tx, err := m.db.BeginTx(ctx, options)
	if err != nil {
		return nil, m.Normalize(err)
	}

	return &Tx{
		tx:        tx,
		normalize: m.Normalize,
	}, nil
}

func (m *Manager) WithTx(
	ctx context.Context,
	options *sql.TxOptions,
	execute func(*Tx) error,
) (err error) {
	tx, err := m.BeginTx(ctx, options)
	if err != nil {
		return err
	}

	rollback := true
	defer func() {
		if recovered := recover(); recovered != nil {
			if rollback {
				_ = tx.Rollback()
			}
			panic(recovered)
		}

		if !rollback {
			return
		}

		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			if err != nil {
				err = errors.Join(err, rollbackErr)
				return
			}
			err = rollbackErr
		}
	}()

	if callbackErr := execute(tx); callbackErr != nil {
		return callbackErr
	}

	rollback = false
	return tx.Commit()
}
