package sql_driver_manager

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

var txTestDriverID atomic.Uint64

type txTestState struct {
	mu          sync.Mutex
	begin       int
	commit      int
	rollback    int
	commitErr   error
	rollbackErr error
	execErr     error
}

func (s *txTestState) recordBegin() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.begin++
}

func (s *txTestState) recordCommit() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.commit++
	return s.commitErr
}

func (s *txTestState) recordRollback() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rollback++
	return s.rollbackErr
}

func (s *txTestState) recordExec() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.execErr
}

func (s *txTestState) counts() (begin, commit, rollback int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.begin, s.commit, s.rollback
}

type txTestDriver struct{ state *txTestState }

func (d txTestDriver) Open(string) (driver.Conn, error) {
	return &txTestConn{state: d.state}, nil
}

type txTestConn struct{ state *txTestState }

func (c *txTestConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepare not supported")
}
func (c *txTestConn) Close() error { return nil }
func (c *txTestConn) Begin() (driver.Tx, error) {
	return c.BeginTx(context.Background(), driver.TxOptions{})
}

func (c *txTestConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	c.state.recordBegin()
	return &txTestTx{state: c.state}, nil
}

func (c *txTestConn) ExecContext(context.Context, string, []driver.NamedValue) (driver.Result, error) {
	if err := c.state.recordExec(); err != nil {
		return nil, err
	}
	return driver.RowsAffected(1), nil
}

type txTestTx struct{ state *txTestState }

func (tx *txTestTx) Commit() error   { return tx.state.recordCommit() }
func (tx *txTestTx) Rollback() error { return tx.state.recordRollback() }

type txTestAdapter struct{ driverName string }

func (a txTestAdapter) Name() string          { return "tx-test" }
func (a txTestAdapter) SQLDriverName() string { return a.driverName }
func (a txTestAdapter) Normalize(err error) *Error {
	if err == nil {
		return nil
	}
	return &Error{
		Kind:    ErrorUnknown,
		Driver:  a.Name(),
		Message: "normalized: " + err.Error(),
		Cause:   err,
	}
}

func newTxTestManager(t *testing.T, state *txTestState) *Manager {
	t.Helper()

	driverName := fmt.Sprintf("gorix_tx_test_%d", txTestDriverID.Add(1))
	sql.Register(driverName, txTestDriver{state: state})

	db, err := sql.Open(driverName, "")
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	return &Manager{
		db:      db,
		adapter: txTestAdapter{driverName: driverName},
	}
}

func TestWithTxReturnsCallbackErrorWithoutNormalizing(t *testing.T) {
	state := &txTestState{}
	manager := newTxTestManager(t, state)
	domainErr := errors.New("domain validation failed")

	err := manager.WithTx(context.Background(), nil, func(*Tx) error {
		return domainErr
	})

	if !errors.Is(err, domainErr) {
		t.Fatalf("expected callback error, got %v", err)
	}
	if _, ok := AsError(err); ok {
		t.Fatalf("expected domain error to remain unnormalized, got %T", err)
	}

	begin, commit, rollback := state.counts()
	if begin != 1 || commit != 0 || rollback != 1 {
		t.Fatalf("unexpected tx counts: begin=%d commit=%d rollback=%d", begin, commit, rollback)
	}
}

func TestWithTxPanicRollsBackAndRepanics(t *testing.T) {
	state := &txTestState{}
	manager := newTxTestManager(t, state)
	panicValue := "callback panic"

	defer func() {
		recovered := recover()
		if recovered != panicValue {
			t.Fatalf("expected panic %q, got %v", panicValue, recovered)
		}

		begin, commit, rollback := state.counts()
		if begin != 1 || commit != 0 || rollback != 1 {
			t.Fatalf("unexpected tx counts: begin=%d commit=%d rollback=%d", begin, commit, rollback)
		}
	}()

	_ = manager.WithTx(context.Background(), nil, func(*Tx) error {
		panic(panicValue)
	})
}

func TestWithTxReturnsCallbackAndRollbackErrors(t *testing.T) {
	rollbackCause := errors.New("rollback failed")
	state := &txTestState{rollbackErr: rollbackCause}
	manager := newTxTestManager(t, state)
	domainErr := errors.New("domain validation failed")

	err := manager.WithTx(context.Background(), nil, func(*Tx) error {
		return domainErr
	})

	if !errors.Is(err, domainErr) {
		t.Fatalf("expected callback error in result, got %v", err)
	}
	if !errors.Is(err, rollbackCause) {
		t.Fatalf("expected rollback cause in result, got %v", err)
	}
	normalizedErr, ok := AsError(err)
	if !ok {
		t.Fatalf("expected normalized rollback error in result, got %T", err)
	}
	if normalizedErr.Cause != rollbackCause {
		t.Fatalf("expected rollback cause %v, got %v", rollbackCause, normalizedErr.Cause)
	}

	begin, commit, rollback := state.counts()
	if begin != 1 || commit != 0 || rollback != 1 {
		t.Fatalf("unexpected tx counts: begin=%d commit=%d rollback=%d", begin, commit, rollback)
	}
}

func TestWithTxReturnsNormalizedDatabaseOperationError(t *testing.T) {
	execCause := errors.New("exec failed")
	state := &txTestState{execErr: execCause}
	manager := newTxTestManager(t, state)

	err := manager.WithTx(context.Background(), nil, func(tx *Tx) error {
		_, err := tx.ExecContext(context.Background(), "insert into users(id) values(?)", 1)
		return err
	})

	if !errors.Is(err, execCause) {
		t.Fatalf("expected exec cause in result, got %v", err)
	}
	normalizedErr, ok := AsError(err)
	if !ok {
		t.Fatalf("expected normalized database operation error, got %T", err)
	}
	if normalizedErr.Cause != execCause {
		t.Fatalf("expected exec cause %v, got %v", execCause, normalizedErr.Cause)
	}

	begin, commit, rollback := state.counts()
	if begin != 1 || commit != 0 || rollback != 1 {
		t.Fatalf("unexpected tx counts: begin=%d commit=%d rollback=%d", begin, commit, rollback)
	}
}

func TestWithTxReturnsNormalizedCommitError(t *testing.T) {
	commitCause := errors.New("commit failed")
	state := &txTestState{commitErr: commitCause}
	manager := newTxTestManager(t, state)

	err := manager.WithTx(context.Background(), nil, func(*Tx) error {
		return nil
	})

	if !errors.Is(err, commitCause) {
		t.Fatalf("expected commit cause in result, got %v", err)
	}
	normalizedErr, ok := AsError(err)
	if !ok {
		t.Fatalf("expected normalized commit error, got %T", err)
	}
	if normalizedErr.Cause != commitCause {
		t.Fatalf("expected commit cause %v, got %v", commitCause, normalizedErr.Cause)
	}

	begin, commit, rollback := state.counts()
	if begin != 1 || commit != 1 || rollback != 0 {
		t.Fatalf("unexpected tx counts: begin=%d commit=%d rollback=%d", begin, commit, rollback)
	}
}
