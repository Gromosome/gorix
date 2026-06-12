package test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"testing"
	"time"

	sqldriver "github.com/Gromosome/gorix/sql-driver-manager"
)

var managerTestErr = errors.New("manager test driver error")

type managerTestAdapter struct{}

func (managerTestAdapter) Name() string          { return "manager-test" }
func (managerTestAdapter) SQLDriverName() string { return "manager-test-sql-driver" }
func (managerTestAdapter) Normalize(err error) *sqldriver.Error {
	if err == nil {
		return nil
	}
	return &sqldriver.Error{
		Kind:    sqldriver.ErrorUnknown,
		Driver:  "manager-test",
		Message: err.Error(),
		Cause:   err,
	}
}

type managerTestDriver struct{}

func (managerTestDriver) Open(string) (driver.Conn, error) {
	return &managerTestConn{}, nil
}

type managerTestConn struct{}

func (managerTestConn) Prepare(query string) (driver.Stmt, error) {
	return managerTestStmt{query: query}, nil
}

func (managerTestConn) Close() error {
	return nil
}

func (managerTestConn) Begin() (driver.Tx, error) {
	return managerTestTx{}, nil
}

func (managerTestConn) Ping(context.Context) error {
	return nil
}

func (managerTestConn) ExecContext(
	_ context.Context,
	query string,
	_ []driver.NamedValue,
) (driver.Result, error) {
	if query == "fail" {
		return nil, managerTestErr
	}
	return managerTestResult(1), nil
}

func (managerTestConn) QueryContext(
	_ context.Context,
	query string,
	_ []driver.NamedValue,
) (driver.Rows, error) {
	if query == "fail" {
		return nil, managerTestErr
	}
	return &managerTestRows{values: []driver.Value{int64(42)}}, nil
}

func (managerTestConn) PrepareContext(
	_ context.Context,
	query string,
) (driver.Stmt, error) {
	if query == "fail" {
		return nil, managerTestErr
	}
	return managerTestStmt{query: query}, nil
}

func (managerTestConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	return managerTestTx{}, nil
}

type managerTestResult int64

func (r managerTestResult) LastInsertId() (int64, error) {
	return int64(r), nil
}

func (r managerTestResult) RowsAffected() (int64, error) {
	return int64(r), nil
}

type managerTestRows struct {
	values []driver.Value
	read   bool
}

func (r *managerTestRows) Columns() []string {
	return []string{"value"}
}

func (r *managerTestRows) Close() error {
	return nil
}

func (r *managerTestRows) Next(dest []driver.Value) error {
	if r.read {
		return io.EOF
	}
	r.read = true
	dest[0] = r.values[0]
	return nil
}

type managerTestStmt struct {
	query string
}

func (s managerTestStmt) Close() error {
	return nil
}

func (s managerTestStmt) NumInput() int {
	return -1
}

func (s managerTestStmt) Exec([]driver.Value) (driver.Result, error) {
	if s.query == "fail" {
		return nil, managerTestErr
	}
	return managerTestResult(1), nil
}

func (s managerTestStmt) Query([]driver.Value) (driver.Rows, error) {
	if s.query == "fail" {
		return nil, managerTestErr
	}
	return &managerTestRows{values: []driver.Value{int64(7)}}, nil
}

type managerTestTx struct{}

func (managerTestTx) Commit() error {
	return nil
}

func (managerTestTx) Rollback() error {
	return nil
}

func init() {
	sql.Register("manager-test-sql-driver", managerTestDriver{})
	sqldriver.Register(managerTestAdapter{})
}

func TestOpenConfiguresAndExposesManager(t *testing.T) {
	manager, err := sqldriver.Open(context.Background(), sqldriver.Config{
		Driver:          " manager-test ",
		DSN:             "ignored",
		MaxOpenConns:    3,
		MaxIdleConns:    2,
		ConnMaxLifetime: time.Minute,
		ConnMaxIdleTime: time.Second,
		PingTimeout:     time.Second,
	})
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}
	defer manager.Close()

	if manager.Driver() != "manager-test" {
		t.Fatalf("unexpected manager driver: %s", manager.Driver())
	}

	if stats := manager.Stats(); stats.MaxOpenConnections != 3 {
		t.Fatalf("unexpected max open connections: %d", stats.MaxOpenConnections)
	}

	if err = manager.PingContext(context.Background()); err != nil {
		t.Fatalf("PingContext returned error: %v", err)
	}
}

func TestManagerOperationsUseNormalizer(t *testing.T) {
	manager, err := sqldriver.Open(context.Background(), sqldriver.Config{
		Driver: "manager-test",
		DSN:    "ignored",
	})
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}
	defer manager.Close()

	if _, err = manager.ExecContext(context.Background(), "fail"); !sqldriver.IsKind(err, sqldriver.ErrorUnknown) {
		t.Fatalf("ExecContext error was not normalized: %v", err)
	}
	if !errors.Is(err, managerTestErr) {
		t.Fatalf("ExecContext normalized error does not wrap cause: %v", err)
	}

	if _, err = manager.QueryContext(context.Background(), "fail"); !sqldriver.IsKind(err, sqldriver.ErrorUnknown) {
		t.Fatalf("QueryContext error was not normalized: %v", err)
	}

	if _, err = manager.PrepareContext(context.Background(), "fail"); !sqldriver.IsKind(err, sqldriver.ErrorUnknown) {
		t.Fatalf("PrepareContext error was not normalized: %v", err)
	}
}

func TestManagerQueryPrepareAndTxSuccessPaths(t *testing.T) {
	manager, err := sqldriver.Open(context.Background(), sqldriver.Config{
		Driver: "manager-test",
		DSN:    "ignored",
	})
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}
	defer manager.Close()

	var value int
	if err = manager.QueryRowContext(context.Background(), "select").Scan(&value); err != nil {
		t.Fatalf("QueryRowContext Scan returned error: %v", err)
	}
	if value != 42 {
		t.Fatalf("unexpected query row value: %d", value)
	}

	statement, err := manager.PrepareContext(context.Background(), "select")
	if err != nil {
		t.Fatalf("PrepareContext returned error: %v", err)
	}
	defer statement.Close()

	rows, err := statement.QueryContext(context.Background())
	if err != nil {
		t.Fatalf("Stmt QueryContext returned error: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		t.Fatal("expected prepared statement row")
	}
	if err = rows.Scan(&value); err != nil {
		t.Fatalf("Rows Scan returned error: %v", err)
	}
	if value != 7 {
		t.Fatalf("unexpected prepared row value: %d", value)
	}

	if err = manager.WithTx(context.Background(), nil, func(tx *sqldriver.Tx) error {
		_, executeErr := tx.ExecContext(context.Background(), "ok")
		return executeErr
	}); err != nil {
		t.Fatalf("WithTx returned error: %v", err)
	}
}
