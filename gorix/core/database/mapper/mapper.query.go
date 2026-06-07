package mapper

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Gromosome/gorix/gorix/core/database"
)

func QueryOne[T any](
	ctx context.Context,
	executor database.Executor,
	query string,
	args ...any,
) (*T, error) {
	if executor == nil {
		return nil, fmt.Errorf(
			"gorix mapper: executor cannot be nil",
		)
	}

	rows, err := executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf(
			"gorix mapper: query failed: %w",
			err,
		)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf(
				"gorix mapper: row iteration failed: %w",
				err,
			)
		}

		return nil, sql.ErrNoRows
	}

	item, err := ScanStruct[T](rows)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func QueryMany[T any](
	ctx context.Context,
	executor database.Executor,
	query string,
	args ...any,
) ([]T, error) {
	if executor == nil {
		return nil, fmt.Errorf(
			"gorix mapper: executor cannot be nil",
		)
	}

	rows, err := executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf(
			"gorix mapper: query failed: %w",
			err,
		)
	}
	defer rows.Close()

	results := make([]T, 0)

	for rows.Next() {
		item, err := ScanStruct[T](rows)
		if err != nil {
			return nil, err
		}

		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"gorix mapper: row iteration failed: %w",
			err,
		)
	}

	return results, nil
}

func Exec(
	ctx context.Context,
	executor database.Executor,
	query string,
	args ...any,
) (sql.Result, error) {
	if executor == nil {
		return nil, fmt.Errorf(
			"gorix mapper: executor cannot be nil",
		)
	}

	result, err := executor.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf(
			"gorix mapper: execution failed: %w",
			err,
		)
	}

	return result, nil
}
