package document_driver_manager

import (
	"context"
)

type Client interface {
	Ping(ctx context.Context) error
	Close(ctx context.Context) error
	Database(name string) Database
}

type Database interface {
	Name() string
	Collection(name string) Collection
}
