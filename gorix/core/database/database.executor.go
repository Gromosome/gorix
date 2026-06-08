package database

import gorixcontext "github.com/Gromosome/gorix/gorix/core/context"

type Executor interface {
	Exec(
		ctx *gorixcontext.Context,
		query string,
		args ...any,
	) Result

	Query(
		ctx *gorixcontext.Context,
		query string,
		args ...any,
	) *Rows

	QueryRow(
		ctx *gorixcontext.Context,
		query string,
		args ...any,
	) *Row
}

var _ Executor = (*DB)(nil)
var _ Executor = (*Tx)(nil)
