package repository

import (
	"fmt"
	"reflect"
	"strings"

	gorixcontext "github.com/Gromosome/gorix/gorix/core/context"
	"github.com/Gromosome/gorix/gorix/core/database"
	"github.com/Gromosome/gorix/gorix/core/database/mapper"
	"github.com/Gromosome/gorix/gorix/core/database/orm"
)

type Repository[T any, ID comparable] struct {
	manager        *database.Manager
	connectionName string
	dialect        orm.Dialect
	metadata       *orm.EntityMetadata
}

func NewRepository[T any, ID comparable](
	manager *database.Manager,
	connectionNames ...string,
) (*Repository[T, ID], error) {
	if manager == nil {
		return nil, fmt.Errorf(
			"gorix orm: database manager cannot be nil",
		)
	}

	connectionName := database.DefaultConnectionName

	if len(connectionNames) > 0 &&
		strings.TrimSpace(connectionNames[0]) != "" {
		connectionName = strings.TrimSpace(
			connectionNames[0],
		)
	}

	connection, err := manager.Connection(
		connectionName,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"gorix orm: failed to resolve connection %q: %w",
			connectionName,
			err,
		)
	}

	dialect, err := orm.ResolveDialect(
		connection.Driver(),
	)
	if err != nil {
		return nil, err
	}

	metadata, err := orm.MetadataOf[T]()
	if err != nil {
		return nil, err
	}

	return &Repository[T, ID]{
		manager:        manager,
		connectionName: connectionName,
		dialect:        dialect,
		metadata:       metadata,
	}, nil
}

func (r *Repository[T, ID]) ConnectionName() string {
	if r == nil {
		return ""
	}

	return r.connectionName
}

func (r *Repository[T, ID]) Metadata() *orm.EntityMetadata {
	if r == nil {
		return nil
	}

	return r.metadata
}

func (r *Repository[T, ID]) DB() (
	*database.DB,
	error,
) {
	if r == nil {
		return nil, fmt.Errorf(
			"gorix orm: repository cannot be nil",
		)
	}

	if r.manager == nil {
		return nil, fmt.Errorf(
			"gorix orm: database manager is unavailable",
		)
	}

	db, err := r.manager.DB(
		r.connectionName,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"gorix orm: failed to resolve database %q: %w",
			r.connectionName,
			err,
		)
	}

	return db, nil
}

func (r *Repository[T, ID]) FindByID(
	ctx *gorixcontext.Context,
	id ID,
) (*T, error) {
	if err := validateRepositoryContext(ctx); err != nil {
		return nil, err
	}

	db, err := r.DB()
	if err != nil {
		return nil, err
	}

	return r.findByIDWithExecutor(
		ctx,
		db,
		id,
	)
}

func (r *Repository[T, ID]) FindAll(
	ctx *gorixcontext.Context,
) ([]T, error) {
	if err := validateRepositoryContext(ctx); err != nil {
		return nil, err
	}

	db, err := r.DB()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(
		"SELECT %s FROM %s",
		r.selectColumns(),
		r.quoteTable(),
	)

	results := make([]T, 0)

	if err := mapper.QueryManyInto(
		ctx,
		db,
		&results,
		query,
	); err != nil {
		return nil, fmt.Errorf(
			"gorix orm: find all failed: %w",
			err,
		)
	}

	return results, nil
}

func (r *Repository[T, ID]) Find(
	ctx *gorixcontext.Context,
	builder *orm.QueryBuilder,
) ([]T, error) {
	if err := validateRepositoryContext(ctx); err != nil {
		return nil, err
	}

	if builder == nil {
		return nil, fmt.Errorf(
			"gorix orm: query builder cannot be nil",
		)
	}

	db, err := r.DB()
	if err != nil {
		return nil, err
	}

	query, args, err := builder.BuildSelect()
	if err != nil {
		return nil, err
	}

	results := make([]T, 0)

	if err := mapper.QueryManyInto(
		ctx,
		db,
		&results,
		query,
		args...,
	); err != nil {
		return nil, fmt.Errorf(
			"gorix orm: find query failed: %w",
			err,
		)
	}

	return results, nil
}

func (r *Repository[T, ID]) FindOne(
	ctx *gorixcontext.Context,
	builder *orm.QueryBuilder,
) (*T, error) {
	if err := validateRepositoryContext(ctx); err != nil {
		return nil, err
	}

	if builder == nil {
		return nil, fmt.Errorf(
			"gorix orm: query builder cannot be nil",
		)
	}

	db, err := r.DB()
	if err != nil {
		return nil, err
	}

	query, args, err := builder.
		Limit(1).
		BuildSelect()
	if err != nil {
		return nil, err
	}

	var entity T

	if err := mapper.QueryOneInto(
		ctx,
		db,
		&entity,
		query,
		args...,
	); err != nil {
		if database.IsNoRows(err) {
			return nil, orm.ErrEntityNotFound
		}

		return nil, fmt.Errorf(
			"gorix orm: find one failed: %w",
			err,
		)
	}

	return &entity, nil
}

func (r *Repository[T, ID]) Insert(
	ctx *gorixcontext.Context,
	entity *T,
) error {
	if err := validateRepositoryContext(ctx); err != nil {
		return err
	}

	db, err := r.DB()
	if err != nil {
		return err
	}

	return r.insertWithExecutor(
		ctx,
		db,
		entity,
	)
}

func (r *Repository[T, ID]) Update(
	ctx *gorixcontext.Context,
	entity *T,
) error {
	if err := validateRepositoryContext(ctx); err != nil {
		return err
	}

	db, err := r.DB()
	if err != nil {
		return err
	}

	return r.updateWithExecutor(
		ctx,
		db,
		entity,
	)
}

func (r *Repository[T, ID]) Save(
	ctx *gorixcontext.Context,
	entity *T,
) error {
	if err := validateRepositoryContext(ctx); err != nil {
		return err
	}

	if entity == nil {
		return fmt.Errorf(
			"gorix orm: entity cannot be nil",
		)
	}

	value, err := entityStructValue(entity)
	if err != nil {
		return err
	}

	primaryKey := r.metadata.PrimaryKey

	if primaryKey == nil {
		return fmt.Errorf(
			"gorix orm: entity %s has no primary key",
			r.metadata.Type.Name(),
		)
	}

	if value.Field(primaryKey.Index).IsZero() {
		return r.Insert(ctx, entity)
	}

	return r.Update(ctx, entity)
}

func (r *Repository[T, ID]) DeleteByID(
	ctx *gorixcontext.Context,
	id ID,
) error {
	if err := validateRepositoryContext(ctx); err != nil {
		return err
	}

	db, err := r.DB()
	if err != nil {
		return err
	}

	return r.deleteByIDWithExecutor(
		ctx,
		db,
		id,
	)
}

func (r *Repository[T, ID]) Delete(
	ctx *gorixcontext.Context,
	entity *T,
) error {
	if entity == nil {
		return fmt.Errorf(
			"gorix orm: entity cannot be nil",
		)
	}

	value, err := entityStructValue(entity)
	if err != nil {
		return err
	}

	primaryKey := r.metadata.PrimaryKey

	if primaryKey == nil {
		return fmt.Errorf(
			"gorix orm: primary key is unavailable",
		)
	}

	primaryKeyValue := value.Field(
		primaryKey.Index,
	)

	if primaryKeyValue.IsZero() {
		return orm.ErrMissingID
	}

	id, ok := primaryKeyValue.Interface().(ID)
	if !ok {
		return fmt.Errorf(
			"gorix orm: primary key type %s does not match repository ID type",
			primaryKeyValue.Type(),
		)
	}

	return r.DeleteByID(ctx, id)
}

func (r *Repository[T, ID]) Count(
	ctx *gorixcontext.Context,
) (int64, error) {
	if err := validateRepositoryContext(ctx); err != nil {
		return 0, err
	}

	db, err := r.DB()
	if err != nil {
		return 0, err
	}

	query := fmt.Sprintf(
		"SELECT COUNT(*) FROM %s",
		r.quoteTable(),
	)

	var count int64

	if err := db.QueryRow(
		ctx,
		query,
	).Scan(&count); err != nil {
		return 0, fmt.Errorf(
			"gorix orm: count failed: %w",
			err,
		)
	}

	return count, nil
}

func (r *Repository[T, ID]) ExistsByID(
	ctx *gorixcontext.Context,
	id ID,
) (bool, error) {
	if err := validateRepositoryContext(ctx); err != nil {
		return false, err
	}

	db, err := r.DB()
	if err != nil {
		return false, err
	}

	primaryKey := r.metadata.PrimaryKey

	query := fmt.Sprintf(
		"SELECT COUNT(*) FROM %s WHERE %s = %s",
		r.quoteTable(),
		r.dialect.QuoteIdentifier(
			primaryKey.ColumnName,
		),
		r.dialect.Placeholder(1),
	)

	var count int64

	if err := db.QueryRow(
		ctx,
		query,
		id,
	).Scan(&count); err != nil {
		return false, fmt.Errorf(
			"gorix orm: exists query failed: %w",
			err,
		)
	}

	return count > 0, nil
}

func (r *Repository[T, ID]) NewQuery() *orm.QueryBuilder {
	if r == nil ||
		r.metadata == nil ||
		r.dialect == nil {
		return nil
	}

	columns := make(
		[]string,
		0,
		len(r.metadata.Fields),
	)

	for _, field := range r.metadata.Fields {
		columns = append(
			columns,
			field.ColumnName,
		)
	}

	return orm.NewQueryBuilder(
		r.dialect,
		r.metadata.TableName,
	).Select(columns...)
}

func (r *Repository[T, ID]) WithExecutor(
	executor database.Executor,
) *ScopedRepository[T, ID] {
	return &ScopedRepository[T, ID]{
		repository: r,
		executor:   executor,
	}
}

func (r *Repository[T, ID]) findByIDWithExecutor(
	ctx *gorixcontext.Context,
	executor database.Executor,
	id ID,
) (*T, error) {
	primaryKey := r.metadata.PrimaryKey

	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s = %s",
		r.selectColumns(),
		r.quoteTable(),
		r.dialect.QuoteIdentifier(
			primaryKey.ColumnName,
		),
		r.dialect.Placeholder(1),
	)

	var entity T

	if err := mapper.QueryOneInto(
		ctx,
		executor,
		&entity,
		query,
		id,
	); err != nil {
		if database.IsNoRows(err) {
			return nil, orm.ErrEntityNotFound
		}

		return nil, fmt.Errorf(
			"gorix orm: find by ID failed: %w",
			err,
		)
	}

	return &entity, nil
}

func (r *Repository[T, ID]) insertWithExecutor(
	ctx *gorixcontext.Context,
	executor database.Executor,
	entity *T,
) error {
	if entity == nil {
		return fmt.Errorf(
			"gorix orm: entity cannot be nil",
		)
	}

	value, err := entityStructValue(entity)
	if err != nil {
		return err
	}

	columns := make([]string, 0)
	placeholders := make([]string, 0)
	args := make([]any, 0)

	for _, field := range r.metadata.Fields {
		if field.ReadOnly ||
			field.AutoIncrement {
			continue
		}

		fieldValue := value.Field(
			field.Index,
		)

		if field.OmitEmpty &&
			fieldValue.IsZero() {
			continue
		}

		columns = append(
			columns,
			r.dialect.QuoteIdentifier(
				field.ColumnName,
			),
		)

		args = append(
			args,
			fieldValue.Interface(),
		)

		placeholders = append(
			placeholders,
			r.dialect.Placeholder(
				len(args),
			),
		)
	}

	if len(columns) == 0 {
		return fmt.Errorf(
			"gorix orm: entity has no insertable fields",
		)
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		r.quoteTable(),
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	primaryKey := r.metadata.PrimaryKey

	if primaryKey != nil &&
		primaryKey.AutoIncrement &&
		r.dialect.SupportsReturning() {
		query += " RETURNING " +
			r.dialect.QuoteIdentifier(
				primaryKey.ColumnName,
			)

		return executor.QueryRow(
			ctx,
			query,
			args...,
		).Scan(
			value.
				Field(primaryKey.Index).
				Addr().
				Interface(),
		)
	}

	result, err := executor.Exec(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return fmt.Errorf(
			"gorix orm: insert failed: %w",
			err,
		)
	}

	if primaryKey != nil &&
		primaryKey.AutoIncrement {
		id, idErr := result.LastInsertID()
		if idErr == nil {
			setIntegerField(
				value.Field(primaryKey.Index),
				id,
			)
		}
	}

	return nil
}

func (r *Repository[T, ID]) updateWithExecutor(
	ctx *gorixcontext.Context,
	executor database.Executor,
	entity *T,
) error {
	if entity == nil {
		return fmt.Errorf(
			"gorix orm: entity cannot be nil",
		)
	}

	value, err := entityStructValue(entity)
	if err != nil {
		return err
	}

	primaryKey := r.metadata.PrimaryKey
	if primaryKey == nil {
		return fmt.Errorf(
			"gorix orm: primary key is unavailable",
		)
	}

	primaryKeyValue := value.Field(
		primaryKey.Index,
	)

	if primaryKeyValue.IsZero() {
		return orm.ErrMissingID
	}

	assignments := make([]string, 0)
	args := make([]any, 0)

	for _, field := range r.metadata.Fields {
		if field.PrimaryKey ||
			field.ReadOnly ||
			field.AutoIncrement {
			continue
		}

		fieldValue := value.Field(
			field.Index,
		)

		if field.OmitEmpty &&
			fieldValue.IsZero() {
			continue
		}

		args = append(
			args,
			fieldValue.Interface(),
		)

		assignments = append(
			assignments,
			fmt.Sprintf(
				"%s = %s",
				r.dialect.QuoteIdentifier(
					field.ColumnName,
				),
				r.dialect.Placeholder(
					len(args),
				),
			),
		)
	}

	if len(assignments) == 0 {
		return nil
	}

	args = append(
		args,
		primaryKeyValue.Interface(),
	)

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s = %s",
		r.quoteTable(),
		strings.Join(assignments, ", "),
		r.dialect.QuoteIdentifier(
			primaryKey.ColumnName,
		),
		r.dialect.Placeholder(
			len(args),
		),
	)

	result, err := executor.Exec(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return fmt.Errorf(
			"gorix orm: update failed: %w",
			err,
		)
	}

	rowsAffected, err := result.RowsAffected()
	if err == nil && rowsAffected == 0 {
		return orm.ErrEntityNotFound
	}

	return nil
}

func (r *Repository[T, ID]) deleteByIDWithExecutor(
	ctx *gorixcontext.Context,
	executor database.Executor,
	id ID,
) error {
	primaryKey := r.metadata.PrimaryKey

	query := fmt.Sprintf(
		"DELETE FROM %s WHERE %s = %s",
		r.quoteTable(),
		r.dialect.QuoteIdentifier(
			primaryKey.ColumnName,
		),
		r.dialect.Placeholder(1),
	)

	result, err := executor.Exec(
		ctx,
		query,
		id,
	)
	if err != nil {
		return fmt.Errorf(
			"gorix orm: delete failed: %w",
			err,
		)
	}

	rowsAffected, err := result.RowsAffected()
	if err == nil && rowsAffected == 0 {
		return orm.ErrEntityNotFound
	}

	return nil
}

func (r *Repository[T, ID]) selectColumns() string {
	columns := make(
		[]string,
		0,
		len(r.metadata.Fields),
	)

	for _, field := range r.metadata.Fields {
		columns = append(
			columns,
			r.dialect.QuoteIdentifier(
				field.ColumnName,
			),
		)
	}

	return strings.Join(
		columns,
		", ",
	)
}

func (r *Repository[T, ID]) quoteTable() string {
	return r.dialect.QuoteIdentifier(
		r.metadata.TableName,
	)
}

type ScopedRepository[T any, ID comparable] struct {
	repository *Repository[T, ID]
	executor   database.Executor
}

func (r *ScopedRepository[T, ID]) FindByID(
	ctx *gorixcontext.Context,
	id ID,
) (*T, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	return r.repository.findByIDWithExecutor(
		ctx,
		r.executor,
		id,
	)
}

func (r *ScopedRepository[T, ID]) Insert(
	ctx *gorixcontext.Context,
	entity *T,
) error {
	if err := r.validate(); err != nil {
		return err
	}

	return r.repository.insertWithExecutor(
		ctx,
		r.executor,
		entity,
	)
}

func (r *ScopedRepository[T, ID]) Update(
	ctx *gorixcontext.Context,
	entity *T,
) error {
	if err := r.validate(); err != nil {
		return err
	}

	return r.repository.updateWithExecutor(
		ctx,
		r.executor,
		entity,
	)
}

func (r *ScopedRepository[T, ID]) Save(
	ctx *gorixcontext.Context,
	entity *T,
) error {
	if err := r.validate(); err != nil {
		return err
	}

	if entity == nil {
		return fmt.Errorf(
			"gorix orm: entity cannot be nil",
		)
	}

	value, err := entityStructValue(entity)
	if err != nil {
		return err
	}

	primaryKey := r.repository.metadata.PrimaryKey

	if value.Field(primaryKey.Index).IsZero() {
		return r.Insert(ctx, entity)
	}

	return r.Update(ctx, entity)
}

func (r *ScopedRepository[T, ID]) DeleteByID(
	ctx *gorixcontext.Context,
	id ID,
) error {
	if err := r.validate(); err != nil {
		return err
	}

	return r.repository.deleteByIDWithExecutor(
		ctx,
		r.executor,
		id,
	)
}

func (r *ScopedRepository[T, ID]) validate() error {
	if r == nil ||
		r.repository == nil {
		return fmt.Errorf(
			"gorix orm: scoped repository is unavailable",
		)
	}

	if r.executor == nil {
		return fmt.Errorf(
			"gorix orm: scoped executor is unavailable",
		)
	}

	return nil
}

func validateRepositoryContext(
	ctx *gorixcontext.Context,
) error {
	if ctx == nil {
		return fmt.Errorf(
			"gorix orm: context cannot be nil",
		)
	}

	if err := ctx.Err(); err != nil {
		return fmt.Errorf(
			"gorix orm: context is closed: %w",
			err,
		)
	}

	return nil
}

func entityStructValue[T any](
	entity *T,
) (reflect.Value, error) {
	if entity == nil {
		return reflect.Value{}, fmt.Errorf(
			"gorix orm: entity cannot be nil",
		)
	}

	value := reflect.ValueOf(entity)

	if value.Kind() != reflect.Pointer ||
		value.IsNil() {
		return reflect.Value{}, fmt.Errorf(
			"gorix orm: entity must be a non-nil pointer",
		)
	}

	value = value.Elem()

	if value.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf(
			"gorix orm: entity must point to a struct, got %s",
			value.Kind(),
		)
	}

	return value, nil
}

func setIntegerField(
	field reflect.Value,
	value int64,
) {
	if !field.IsValid() ||
		!field.CanSet() {
		return
	}

	switch field.Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		field.SetInt(value)

	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		field.SetUint(uint64(value))
	}
}
