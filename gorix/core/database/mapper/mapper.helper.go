package mapper

import (
	"context"

	"github.com/Gromosome/gorix/gorix/core/database"
)

func QueryNamedOne[T any](
	ctx context.Context,
	executor database.Executor,
	registry *StatementRegistry,
	name string,
	args ...any,
) (*T, error) {
	query, err := registry.Get(name)
	if err != nil {
		return nil, err
	}

	return QueryOne[T](ctx, executor, query, args...)
}

func QueryNamedMany[T any](
	ctx context.Context,
	executor database.Executor,
	registry *StatementRegistry,
	name string,
	args ...any,
) ([]T, error) {
	query, err := registry.Get(name)
	if err != nil {
		return nil, err
	}

	return QueryMany[T](ctx, executor, query, args...)
}
