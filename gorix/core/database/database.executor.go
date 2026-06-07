package database

import gorixcontext "github.com/Gromosome/gorix/gorix/core/context"

type Executor interface {
	Exec(
		ctx *gorixcontext.Context,
		query string,
		args ...any,
	) (Result, error)

	Query(
		ctx *gorixcontext.Context,
		query string,
		args ...any,
	) (*Rows, error)

	QueryRow(
		ctx *gorixcontext.Context,
		query string,
		args ...any,
	) *Row
}
