package database

import (
	"database/sql"
)

type DB struct {
	native *sql.DB
}

type Tx struct {
	native *sql.Tx
}

type Rows struct {
	native *sql.Rows
	err    error
}

type Row struct {
	native *sql.Row
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
