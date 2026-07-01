package mongo

import (
	"context"
	"fmt"
	"strings"

	docdriver "github.com/Gromosome/gorix/document-driver-manager"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (d *Database) EnsureSchema(
	ctx context.Context,
	schema docdriver.Schema,
) error {
	if d == nil || d.client == nil || d.client.native == nil {
		return fmt.Errorf("gorix mongo: database is unavailable")
	}

	collectionName := normalizeCollectionName(schema.Collection)
	if collectionName == "" {
		return fmt.Errorf("gorix mongo: collection name is required")
	}

	nativeDB := d.client.native.Database(d.name)

	names, err := nativeDB.ListCollectionNames(
		ctx,
		bson.M{"name": collectionName},
	)
	if err != nil {
		return d.client.adapter.Normalize(err)
	}

	if len(names) == 0 {
		if err := nativeDB.CreateCollection(ctx, collectionName); err != nil {
			if !mongodriver.IsDuplicateKeyError(err) {
				return d.client.adapter.Normalize(err)
			}
		}
	}

	if len(schema.Indexes) == 0 {
		return nil
	}

	models := make(
		[]mongodriver.IndexModel,
		0,
		len(schema.Indexes),
	)

	for _, index := range schema.Indexes {
		keys := bson.D{}

		for _, field := range index.Fields {
			fieldName := strings.TrimSpace(field.Field)
			if fieldName == "" || fieldName == "_id" || fieldName == "id" {
				continue
			}

			direction := int32(1)
			if field.Desc {
				direction = -1
			}

			keys = append(
				keys,
				bson.E{
					Key:   fieldName,
					Value: direction,
				},
			)
		}

		if len(keys) == 0 {
			continue
		}

		name := strings.TrimSpace(index.Name)
		if name == "" {
			name = keys[0].Key
		}

		models = append(
			models,
			mongodriver.IndexModel{
				Keys: keys,
				Options: options.Index().
					SetName(name).
					SetUnique(index.Unique),
			},
		)
	}

	if len(models) == 0 {
		return nil
	}

	_, err = nativeDB.
		Collection(collectionName).
		Indexes().
		CreateMany(ctx, models)

	if err != nil {
		return d.client.adapter.Normalize(err)
	}

	return nil
}
