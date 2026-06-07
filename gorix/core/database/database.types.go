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
}

type Row struct {
	native *sql.Row
}

type Result struct {
	native sql.Result
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
