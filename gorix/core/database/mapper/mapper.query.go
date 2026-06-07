package mapper

import (
	"fmt"
	"reflect"
	"strings"

	gorixcontext "github.com/Gromosome/gorix/gorix/core/context"
	"github.com/Gromosome/gorix/gorix/core/database"
)

// QueryOneInto executes a query and maps the first returned row into target.
//
// Target must be a non-nil pointer to a struct.
func QueryOneInto(
	ctx *gorixcontext.Context,
	executor database.Executor,
	target any,
	query string,
	args ...any,
) error {
	if err := validateMapperContext(ctx); err != nil {
		return err
	}

	if err := validateExecutor(executor); err != nil {
		return err
	}

	if err := validateQuery(query); err != nil {
		return err
	}

	if err := validateSingleTarget(target); err != nil {
		return err
	}

	rows, err := executor.Query(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return fmt.Errorf(
			"gorix mapper: query execution failed: %w",
			err,
		)
	}
	defer func() {
		_ = rows.Close()
	}()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return fmt.Errorf(
				"gorix mapper: row iteration failed: %w",
				err,
			)
		}

		return database.ErrNoRows
	}

	if err := ScanInto(rows, target); err != nil {
		return fmt.Errorf(
			"gorix mapper: row mapping failed: %w",
			err,
		)
	}

	return nil
}

// QueryManyInto executes a query and maps all returned rows into target.
//
// Target must be a non-nil pointer to []T or []*T.
func QueryManyInto(
	ctx *gorixcontext.Context,
	executor database.Executor,
	target any,
	query string,
	args ...any,
) error {
	if err := validateMapperContext(ctx); err != nil {
		return err
	}

	if err := validateExecutor(executor); err != nil {
		return err
	}

	if err := validateQuery(query); err != nil {
		return err
	}

	targetValue, sliceValue, structType, pointerElements, err :=
		resolveSliceTarget(target)
	if err != nil {
		return err
	}

	rows, err := executor.Query(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return fmt.Errorf(
			"gorix mapper: query execution failed: %w",
			err,
		)
	}
	defer func() {
		_ = rows.Close()
	}()

	result := reflect.MakeSlice(
		sliceValue.Type(),
		0,
		0,
	)

	for rows.Next() {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf(
				"gorix mapper: query context ended: %w",
				err,
			)
		}

		itemPointer := reflect.New(structType)

		if err := ScanInto(
			rows,
			itemPointer.Interface(),
		); err != nil {
			return fmt.Errorf(
				"gorix mapper: row mapping failed: %w",
				err,
			)
		}

		if pointerElements {
			result = reflect.Append(
				result,
				itemPointer,
			)
		} else {
			result = reflect.Append(
				result,
				itemPointer.Elem(),
			)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf(
			"gorix mapper: row iteration failed: %w",
			err,
		)
	}

	targetValue.Elem().Set(result)

	return nil
}

// Exec executes a non-query SQL operation through a Gorix executor.
func Exec(
	ctx *gorixcontext.Context,
	executor database.Executor,
	query string,
	args ...any,
) (database.Result, error) {
	if err := validateMapperContext(ctx); err != nil {
		return database.Result{}, err
	}

	if err := validateExecutor(executor); err != nil {
		return database.Result{}, err
	}

	if err := validateQuery(query); err != nil {
		return database.Result{}, err
	}

	result, err := executor.Exec(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return database.Result{}, fmt.Errorf(
			"gorix mapper: statement execution failed: %w",
			err,
		)
	}

	return result, nil
}

// QueryNamedOneInto resolves a statement from the registry and maps one row.
func QueryNamedOneInto(
	ctx *gorixcontext.Context,
	executor database.Executor,
	registry *StatementRegistry,
	target any,
	statementName string,
	args ...any,
) error {
	if registry == nil {
		return fmt.Errorf(
			"gorix mapper: statement registry cannot be nil",
		)
	}

	query, err := registry.Get(statementName)
	if err != nil {
		return err
	}

	return QueryOneInto(
		ctx,
		executor,
		target,
		query,
		args...,
	)
}

// QueryNamedManyInto resolves a statement from the registry and maps all rows.
func QueryNamedManyInto(
	ctx *gorixcontext.Context,
	executor database.Executor,
	registry *StatementRegistry,
	target any,
	statementName string,
	args ...any,
) error {
	if registry == nil {
		return fmt.Errorf(
			"gorix mapper: statement registry cannot be nil",
		)
	}

	query, err := registry.Get(statementName)
	if err != nil {
		return err
	}

	return QueryManyInto(
		ctx,
		executor,
		target,
		query,
		args...,
	)
}

// ExecNamed resolves and executes a statement from the supplied registry.
func ExecNamed(
	ctx *gorixcontext.Context,
	executor database.Executor,
	registry *StatementRegistry,
	statementName string,
	args ...any,
) (database.Result, error) {
	if registry == nil {
		return database.Result{}, fmt.Errorf(
			"gorix mapper: statement registry cannot be nil",
		)
	}

	query, err := registry.Get(statementName)
	if err != nil {
		return database.Result{}, err
	}

	return Exec(
		ctx,
		executor,
		query,
		args...,
	)
}

func validateMapperContext(
	ctx *gorixcontext.Context,
) error {
	if ctx == nil {
		return fmt.Errorf(
			"gorix mapper: context cannot be nil",
		)
	}

	if err := ctx.Err(); err != nil {
		return fmt.Errorf(
			"gorix mapper: context is already closed: %w",
			err,
		)
	}

	return nil
}

func validateExecutor(
	executor database.Executor,
) error {
	if executor == nil {
		return fmt.Errorf(
			"gorix mapper: executor cannot be nil",
		)
	}

	value := reflect.ValueOf(executor)

	switch value.Kind() {
	case reflect.Pointer,
		reflect.Interface,
		reflect.Map,
		reflect.Slice,
		reflect.Func,
		reflect.Chan:
		if value.IsNil() {
			return fmt.Errorf(
				"gorix mapper: executor cannot be nil",
			)
		}
	}

	return nil
}

func validateQuery(query string) error {
	if strings.TrimSpace(query) == "" {
		return fmt.Errorf(
			"gorix mapper: query cannot be empty",
		)
	}

	return nil
}

func validateSingleTarget(
	target any,
) error {
	if target == nil {
		return fmt.Errorf(
			"gorix mapper: QueryOne target cannot be nil",
		)
	}

	targetValue := reflect.ValueOf(target)

	if targetValue.Kind() != reflect.Pointer ||
		targetValue.IsNil() {
		return fmt.Errorf(
			"gorix mapper: QueryOne target must be a non-nil pointer",
		)
	}

	value := targetValue.Elem()

	if value.Kind() != reflect.Struct {
		return fmt.Errorf(
			"gorix mapper: QueryOne target must point to a struct, got %s",
			value.Kind(),
		)
	}

	return nil
}

func resolveSliceTarget(
	target any,
) (
	targetValue reflect.Value,
	sliceValue reflect.Value,
	structType reflect.Type,
	pointerElements bool,
	err error,
) {
	if target == nil {
		err = fmt.Errorf(
			"gorix mapper: QueryMany target cannot be nil",
		)
		return
	}

	targetValue = reflect.ValueOf(target)

	if targetValue.Kind() != reflect.Pointer ||
		targetValue.IsNil() {
		err = fmt.Errorf(
			"gorix mapper: QueryMany target must be a non-nil pointer",
		)
		return
	}

	sliceValue = targetValue.Elem()

	if sliceValue.Kind() != reflect.Slice {
		err = fmt.Errorf(
			"gorix mapper: QueryMany target must point to a slice, got %s",
			sliceValue.Kind(),
		)
		return
	}

	elementType := sliceValue.Type().Elem()

	pointerElements = elementType.Kind() == reflect.Pointer

	structType = elementType

	if pointerElements {
		structType = elementType.Elem()
	}

	if structType.Kind() != reflect.Struct {
		err = fmt.Errorf(
			"gorix mapper: QueryMany slice elements must be structs or pointers to structs, got %s",
			structType.Kind(),
		)
		return
	}

	return
}
