package couchdb

import (
	"context"
	"fmt"
	"strings"

	docdriver "github.com/Gromosome/gorix/document-driver-manager"
)

func (d *Database) EnsureSchema(
	ctx context.Context,
	schema docdriver.Schema,
) error {
	if d == nil || d.client == nil || d.client.native == nil {
		return fmt.Errorf("gorix couchdb: database is unavailable")
	}

	collectionName := normalizeDatabaseName(schema.Collection)
	if collectionName == "" {
		return fmt.Errorf("gorix couchdb: collection name is required")
	}

	dbName := d.collectionDatabaseName(collectionName)

	exists, err := d.client.native.DBExists(ctx, dbName)
	if err != nil {
		return d.client.adapter.Normalize(err)
	}

	if !exists {
		if err := d.client.native.CreateDB(ctx, dbName); err != nil {
			return d.client.adapter.Normalize(err)
		}
	}

	db := d.client.native.DB(dbName)
	if err := db.Err(); err != nil {
		return d.client.adapter.Normalize(err)
	}

	for _, index := range schema.Indexes {
		fields := make([]map[string]string, 0, len(index.Fields))

		for _, field := range index.Fields {
			fieldName := strings.TrimSpace(field.Field)
			if fieldName == "" || fieldName == "_id" || fieldName == "id" {
				continue
			}

			direction := "asc"
			if field.Desc {
				direction = "desc"
			}

			fields = append(
				fields,
				map[string]string{
					fieldName: direction,
				},
			)
		}

		if len(fields) == 0 {
			continue
		}

		name := strings.TrimSpace(index.Name)
		if name == "" {
			name = fields[0][index.Fields[0].Field]
		}

		indexBody := map[string]any{
			"fields": fields,
		}

		if err := db.CreateIndex(
			ctx,
			"gorix",
			name,
			indexBody,
		); err != nil {
			return d.client.adapter.Normalize(err)
		}
	}

	return nil
}
