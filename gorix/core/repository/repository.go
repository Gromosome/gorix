package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

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
		connectionNames[0] != "" {
		connectionName = connectionNames[0]
	}

	connection, err := manager.Connection(connectionName)
	if err != nil {
		return nil, err
	}

	dialect, err := orm.ResolveDialect(connection.Driver())
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

func (r *Repository[T, ID]) DB() (*sql.DB, error) {
	return r.manager.DB(r.connectionName)
}

func (r *Repository[T, ID]) FindByID(
	ctx context.Context,
	id ID,
) (*T, error) {
	db, err := r.DB()
	if err != nil {
		return nil, err
	}

	primaryKey := r.metadata.PrimaryKey

	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s = %s",
		r.selectColumns(),
		r.dialect.QuoteIdentifier(r.metadata.TableName),
		r.dialect.QuoteIdentifier(primaryKey.ColumnName),
		r.dialect.Placeholder(1),
	)

	entity, err := mapper.QueryOne[T](
		ctx,
		db,
		query,
		id,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, orm.ErrEntityNotFound
	}

	return entity, err
}

func (r *Repository[T, ID]) FindAll(
	ctx context.Context,
) ([]T, error) {
	db, err := r.DB()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(
		"SELECT %s FROM %s",
		r.selectColumns(),
		r.dialect.QuoteIdentifier(r.metadata.TableName),
	)

	return mapper.QueryMany[T](
		ctx,
		db,
		query,
	)
}

func (r *Repository[T, ID]) Find(
	ctx context.Context,
	builder *orm.QueryBuilder,
) ([]T, error) {
	db, err := r.DB()
	if err != nil {
		return nil, err
	}

	if builder == nil {
		return nil, fmt.Errorf(
			"gorix orm: query builder cannot be nil",
		)
	}

	query, args, err := builder.BuildSelect()
	if err != nil {
		return nil, err
	}

	return mapper.QueryMany[T](
		ctx,
		db,
		query,
		args...,
	)
}

func (r *Repository[T, ID]) NewQuery() *orm.QueryBuilder {
	columns := make([]string, 0, len(r.metadata.Fields))

	for _, field := range r.metadata.Fields {
		columns = append(columns, field.ColumnName)
	}

	return orm.NewQueryBuilder(
		r.dialect,
		r.metadata.TableName,
	).Select(columns...)
}

func (r *Repository[T, ID]) Insert(
	ctx context.Context,
	entity *T,
) error {
	if entity == nil {
		return fmt.Errorf(
			"gorix orm: entity cannot be nil",
		)
	}

	db, err := r.DB()
	if err != nil {
		return err
	}

	value := reflect.ValueOf(entity).Elem()

	columns := make([]string, 0)
	placeholders := make([]string, 0)
	args := make([]any, 0)

	for _, field := range r.metadata.Fields {
		if field.ReadOnly || field.AutoIncrement {
			continue
		}

		fieldValue := value.Field(field.Index)

		if field.OmitEmpty && fieldValue.IsZero() {
			continue
		}

		columns = append(
			columns,
			r.dialect.QuoteIdentifier(field.ColumnName),
		)

		args = append(args, fieldValue.Interface())

		placeholders = append(
			placeholders,
			r.dialect.Placeholder(len(args)),
		)
	}

	if len(columns) == 0 {
		return fmt.Errorf(
			"gorix orm: entity has no insertable fields",
		)
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		r.dialect.QuoteIdentifier(r.metadata.TableName),
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	primaryKey := r.metadata.PrimaryKey

	if primaryKey.AutoIncrement &&
		r.dialect.SupportsReturning() {
		query += " RETURNING " +
			r.dialect.QuoteIdentifier(primaryKey.ColumnName)

		return db.QueryRowContext(
			ctx,
			query,
			args...,
		).Scan(
			value.Field(primaryKey.Index).Addr().Interface(),
		)
	}

	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf(
			"gorix orm: insert failed: %w",
			err,
		)
	}

	if primaryKey.AutoIncrement {
		id, idErr := result.LastInsertId()
		if idErr == nil {
			setIntegerValue(
				value.Field(primaryKey.Index),
				id,
			)
		}
	}

	return nil
}

func (r *Repository[T, ID]) Update(
	ctx context.Context,
	entity *T,
) error {
	if entity == nil {
		return fmt.Errorf(
			"gorix orm: entity cannot be nil",
		)
	}

	db, err := r.DB()
	if err != nil {
		return err
	}

	value := reflect.ValueOf(entity).Elem()
	primaryKey := r.metadata.PrimaryKey
	primaryKeyValue := value.Field(primaryKey.Index)

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

		fieldValue := value.Field(field.Index)

		if field.OmitEmpty && fieldValue.IsZero() {
			continue
		}

		args = append(args, fieldValue.Interface())

		assignments = append(
			assignments,
			fmt.Sprintf(
				"%s = %s",
				r.dialect.QuoteIdentifier(field.ColumnName),
				r.dialect.Placeholder(len(args)),
			),
		)
	}

	if len(assignments) == 0 {
		return nil
	}

	args = append(args, primaryKeyValue.Interface())

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s = %s",
		r.dialect.QuoteIdentifier(r.metadata.TableName),
		strings.Join(assignments, ", "),
		r.dialect.QuoteIdentifier(primaryKey.ColumnName),
		r.dialect.Placeholder(len(args)),
	)

	result, err := db.ExecContext(ctx, query, args...)
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

func (r *Repository[T, ID]) Save(
	ctx context.Context,
	entity *T,
) error {
	if entity == nil {
		return fmt.Errorf(
			"gorix orm: entity cannot be nil",
		)
	}

	value := reflect.ValueOf(entity).Elem()
	primaryKey := r.metadata.PrimaryKey

	if value.Field(primaryKey.Index).IsZero() {
		return r.Insert(ctx, entity)
	}

	return r.Update(ctx, entity)
}

func (r *Repository[T, ID]) DeleteByID(
	ctx context.Context,
	id ID,
) error {
	db, err := r.DB()
	if err != nil {
		return err
	}

	primaryKey := r.metadata.PrimaryKey

	query := fmt.Sprintf(
		"DELETE FROM %s WHERE %s = %s",
		r.dialect.QuoteIdentifier(r.metadata.TableName),
		r.dialect.QuoteIdentifier(primaryKey.ColumnName),
		r.dialect.Placeholder(1),
	)

	result, err := db.ExecContext(ctx, query, id)
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

func (r *Repository[T, ID]) Count(
	ctx context.Context,
) (int64, error) {
	db, err := r.DB()
	if err != nil {
		return 0, err
	}

	query := fmt.Sprintf(
		"SELECT COUNT(*) FROM %s",
		r.dialect.QuoteIdentifier(r.metadata.TableName),
	)

	var count int64

	if err := db.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return 0, fmt.Errorf(
			"gorix orm: count failed: %w",
			err,
		)
	}

	return count, nil
}

func (r *Repository[T, ID]) ExistsByID(
	ctx context.Context,
	id ID,
) (bool, error) {
	db, err := r.DB()
	if err != nil {
		return false, err
	}

	primaryKey := r.metadata.PrimaryKey

	query := fmt.Sprintf(
		"SELECT COUNT(*) FROM %s WHERE %s = %s",
		r.dialect.QuoteIdentifier(r.metadata.TableName),
		r.dialect.QuoteIdentifier(primaryKey.ColumnName),
		r.dialect.Placeholder(1),
	)

	var count int64

	if err := db.QueryRowContext(ctx, query, id).Scan(&count); err != nil {
		return false, fmt.Errorf(
			"gorix orm: exists query failed: %w",
			err,
		)
	}

	return count > 0, nil
}

func (r *Repository[T, ID]) selectColumns() string {
	columns := make([]string, 0, len(r.metadata.Fields))

	for _, field := range r.metadata.Fields {
		columns = append(
			columns,
			r.dialect.QuoteIdentifier(field.ColumnName),
		)
	}

	return strings.Join(columns, ", ")
}

func setIntegerValue(
	field reflect.Value,
	value int64,
) {
	if !field.CanSet() {
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
