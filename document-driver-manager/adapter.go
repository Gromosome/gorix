package document_driver_manager

import "context"

type Adapter interface {
	Name() string
	Open(ctx context.Context, config Config) (Client, error)
	Normalize(error) *Error
}
