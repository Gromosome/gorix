package database

import (
	"database/sql"

	sqldriver "github.com/Gromosome/gorix/sql-driver-manager"
)

type DB struct {
	native *sqldriver.Manager
}

type Tx struct {
	native *sqldriver.Tx
}

type Rows struct {
	native *sqldriver.Rows
	err    error
}

type Row struct {
	native *sqldriver.Row
	err    error
}

type Result struct {
	native sql.Result
	err    error
}

type TxOptions struct {
	Isolation IsolationLevel
	ReadOnly  bool
}

type IsolationLevel int

const (
	IsolationDefault IsolationLevel = iota
	IsolationReadUncommitted
	IsolationReadCommitted
	IsolationWriteCommitted
	IsolationRepeatableRead
	IsolationSnapshot
	IsolationSerializable
	IsolationLinearizable
)
