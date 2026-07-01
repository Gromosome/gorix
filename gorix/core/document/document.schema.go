package document

import (
	"fmt"
	"strings"

	docdriver "github.com/Gromosome/gorix/document-driver-manager"
	gorixcontext "github.com/Gromosome/gorix/gorix/core/context"
)

func (m *DocumentMetadata) DriverSchema(
	database string,
) docdriver.Schema {
	if m == nil {
		return docdriver.Schema{}
	}

	indexes := make(
		[]docdriver.Index,
		0,
		len(m.Indexes),
	)

	for _, index := range m.Indexes {
		fields := make(
			[]docdriver.IndexField,
			0,
			len(index.Fields),
		)

		for _, field := range index.Fields {
			field = strings.TrimSpace(field)
			if field == "" || field == "_id" || field == "id" {
				continue
			}

			fields = append(
				fields,
				docdriver.IndexField{
					Field: field,
					Desc:  false,
				},
			)
		}

		if len(fields) == 0 {
			continue
		}

		name := strings.TrimSpace(index.Name)
		if name == "" {
			name = fields[0].Field
		}

		indexes = append(
			indexes,
			docdriver.Index{
				Name:   name,
				Fields: fields,
				Unique: index.Unique,
			},
		)
	}

	return docdriver.Schema{
		Database:   database,
		Collection: m.CollectionName,
		Indexes:    indexes,
	}
}

func (r *Repository[T, ID]) EnsureSchema(
	ctx *gorixcontext.Context,
) error {
	if r == nil {
		return fmt.Errorf("gorix document: repository cannot be nil")
	}

	if r.metadata == nil {
		return fmt.Errorf("gorix document: repository metadata is missing")
	}

	connection, err := r.manager.Connection(r.connectionName)
	if err != nil {
		return err
	}

	database := connection.Database()
	if database == nil {
		return fmt.Errorf(
			"gorix document: database is unavailable for connection %q",
			r.connectionName,
		)
	}

	schemaManager, ok := database.(docdriver.SchemaManager)
	if !ok {
		return &docdriver.Error{
			Kind:     docdriver.ErrorSchemaUnsupported,
			Driver:   "gorix-document",
			Database: connection.Config().Database,
			Message: fmt.Sprintf(
				"gorix document: schema is not supported for connection %q",
				connection.Name(),
			),
		}
	}

	return schemaManager.EnsureSchema(
		requestContext(ctx),
		r.metadata.DriverSchema(connection.Config().Database),
	)
}
