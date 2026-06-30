package document

import (
	"context"
	"fmt"
	"reflect"

	docdriver "github.com/Gromosome/gorix/document-driver-manager"
	"github.com/Gromosome/gorix/gorix/config"
	gorixcontext "github.com/Gromosome/gorix/gorix/core/context"
)

type Repository[T any, ID comparable] struct {
	manager        *Manager
	connectionName string
	metadata       *DocumentMetadata
	executor       docdriver.Executor
}

func NewRepository[T any, ID comparable](
	manager *Manager,
	connectionNames ...string,
) (*Repository[T, ID], error) {
	if manager == nil {
		return nil, fmt.Errorf("gorix document: manager cannot be nil")
	}

	connectionName := config.DefaultConnectionName
	if len(connectionNames) > 0 && connectionNames[0] != "" {
		connectionName = connectionNames[0]
	}

	metadata, err := ParseMetadata[T]()
	if err != nil {
		return nil, err
	}

	return &Repository[T, ID]{
		manager:        manager,
		connectionName: connectionName,
		metadata:       metadata,
	}, nil
}

func (r *Repository[T, ID]) Metadata() *DocumentMetadata {
	if r == nil {
		return nil
	}
	return r.metadata
}

func (r *Repository[T, ID]) FindByID(
	ctx *gorixcontext.Context,
	id ID,
) (*T, error) {
	collection, err := r.collection()
	if err != nil {
		return nil, err
	}

	out := new(T)

	if err := collection.FindByID(requestContext(ctx), id, out); err != nil {
		return nil, err
	}

	return out, nil
}

func (r *Repository[T, ID]) Find(
	ctx *gorixcontext.Context,
	query *docdriver.Query,
) ([]T, error) {
	collection, err := r.collection()
	if err != nil {
		return nil, err
	}

	filter, options, err := buildQuery(query)
	if err != nil {
		return nil, err
	}

	out := make([]T, 0)

	if err := collection.Find(
		requestContext(ctx),
		filter,
		&out,
		options,
	); err != nil {
		return nil, err
	}

	return out, nil
}

func (r *Repository[T, ID]) Insert(
	ctx *gorixcontext.Context,
	entity *T,
) error {
	if entity == nil {
		return fmt.Errorf("gorix document: entity cannot be nil")
	}

	collection, err := r.collection()
	if err != nil {
		return err
	}

	result, err := collection.InsertOne(requestContext(ctx), entity)
	if err != nil {
		return err
	}

	if result.ID != nil {
		if err := r.metadata.SetID(entity, result.ID); err != nil {
			return err
		}
	}

	if result.Rev != "" {
		if err := r.metadata.SetRev(entity, result.Rev); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository[T, ID]) Update(
	ctx *gorixcontext.Context,
	entity *T,
) error {
	id, err := r.idFromEntity(entity)
	if err != nil {
		return err
	}

	return r.UpdateByID(ctx, id, entity)
}

func (r *Repository[T, ID]) UpdateByID(
	ctx *gorixcontext.Context,
	id ID,
	entity *T,
) error {
	if entity == nil {
		return fmt.Errorf("gorix document: entity cannot be nil")
	}

	collection, err := r.collection()
	if err != nil {
		return err
	}

	if err := r.metadata.SetID(entity, id); err != nil {
		return err
	}

	result, err := collection.UpdateByID(
		requestContext(ctx),
		id,
		r.metadata.RevValue(entity),
		entity,
	)
	if err != nil {
		return err
	}

	if result.Rev != "" {
		if err := r.metadata.SetRev(entity, result.Rev); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository[T, ID]) Save(
	ctx *gorixcontext.Context,
	entity *T,
) error {
	if entity == nil {
		return fmt.Errorf("gorix document: entity cannot be nil")
	}

	if r.metadata.IsIDZero(entity) {
		return r.Insert(ctx, entity)
	}

	return r.Update(ctx, entity)
}

func (r *Repository[T, ID]) Delete(
	ctx *gorixcontext.Context,
	entity *T,
) error {
	id, err := r.idFromEntity(entity)
	if err != nil {
		return err
	}

	collection, err := r.collection()
	if err != nil {
		return err
	}

	_, err = collection.DeleteByID(
		requestContext(ctx),
		id,
		r.metadata.RevValue(entity),
	)

	return err
}

func (r *Repository[T, ID]) DeleteByID(
	ctx *gorixcontext.Context,
	id ID,
) error {
	collection, err := r.collection()
	if err != nil {
		return err
	}

	_, err = collection.DeleteByID(requestContext(ctx), id, "")
	return err
}

func (r *Repository[T, ID]) Count(
	ctx *gorixcontext.Context,
	query *docdriver.Query,
) (int64, error) {
	collection, err := r.collection()
	if err != nil {
		return 0, err
	}

	filter, _, err := buildQuery(query)
	if err != nil {
		return 0, err
	}

	return collection.Count(requestContext(ctx), filter)
}

func (r *Repository[T, ID]) collection() (docdriver.Collection, error) {
	if r == nil {
		return nil, fmt.Errorf("gorix document: repository cannot be nil")
	}

	if r.metadata == nil {
		return nil, fmt.Errorf("gorix document: repository metadata is missing")
	}

	if r.executor != nil {
		collection := r.executor.Collection(r.metadata.CollectionName)
		if collection == nil {
			return nil, fmt.Errorf(
				"gorix document: collection %q is unavailable",
				r.metadata.CollectionName,
			)
		}
		return collection, nil
	}

	connection, err := r.manager.Connection(r.connectionName)
	if err != nil {
		return nil, err
	}

	collection := connection.Collection(r.metadata.CollectionName)
	if collection == nil {
		return nil, fmt.Errorf(
			"gorix document: collection %q is unavailable",
			r.metadata.CollectionName,
		)
	}

	return collection, nil
}

func (r *Repository[T, ID]) idFromEntity(entity *T) (ID, error) {
	var zero ID

	if entity == nil {
		return zero, fmt.Errorf("gorix document: entity cannot be nil")
	}

	if r.metadata.IsIDZero(entity) {
		return zero, fmt.Errorf("gorix document: entity id cannot be empty")
	}

	idValue, err := r.metadata.IDValue(entity)
	if err != nil {
		return zero, err
	}

	id, ok := idValue.(ID)
	if ok {
		return id, nil
	}

	source := reflect.ValueOf(idValue)
	targetType := reflect.TypeOf(zero)

	if source.IsValid() &&
		targetType != nil &&
		source.Type().ConvertibleTo(targetType) {
		return source.Convert(targetType).Interface().(ID), nil
	}

	return zero, fmt.Errorf(
		"gorix document: cannot convert entity id %T to repository id type",
		idValue,
	)
}

func buildQuery(
	query *docdriver.Query,
) (docdriver.Filter, docdriver.FindOptions, error) {
	if query == nil {
		return docdriver.Filter{}, docdriver.FindOptions{}, nil
	}

	return query.Build()
}

func requestContext(ctx *gorixcontext.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}

	return ctx
}
