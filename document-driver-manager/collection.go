package document_driver_manager

import "context"

type Collection interface {
	Name() string

	InsertOne(
		ctx context.Context,
		document any,
	) (InsertResult, error)

	FindByID(
		ctx context.Context,
		id any,
		out any,
	) error

	Find(
		ctx context.Context,
		filter Filter,
		out any,
		options FindOptions,
	) error

	UpdateByID(
		ctx context.Context,
		id any,
		rev string,
		document any,
	) (UpdateResult, error)

	DeleteByID(
		ctx context.Context,
		id any,
		rev string,
	) (DeleteResult, error)

	Count(
		ctx context.Context,
		filter Filter,
	) (int64, error)
}
