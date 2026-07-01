package document_driver_manager

import "context"

const (
	ErrorSchemaUnsupported ErrorKind = "schema_unsupported"
)

type IndexField struct {
	Field string
	Desc  bool
}

type Index struct {
	Name   string
	Fields []IndexField
	Unique bool
}

type Schema struct {
	Database   string
	Collection string
	Indexes    []Index
}

type SchemaManager interface {
	EnsureSchema(ctx context.Context, schema Schema) error
}
