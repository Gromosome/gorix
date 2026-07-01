package mongo

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	docdriver "github.com/Gromosome/gorix/document-driver-manager"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Collection struct {
	name    string
	native  *mongodriver.Collection
	adapter Adapter
}

func (c *Collection) Name() string {
	if c == nil {
		return ""
	}

	return c.name
}

func (c *Collection) InsertOne(
	ctx context.Context,
	document any,
) (docdriver.InsertResult, error) {
	if err := c.validate(document); err != nil {
		return docdriver.InsertResult{}, err
	}

	result, err := c.native.InsertOne(ctx, document)
	if err != nil {
		return docdriver.InsertResult{},
			c.adapter.Normalize(err)
	}

	setDocumentID(document, result.InsertedID)

	return docdriver.InsertResult{
		ID: result.InsertedID,
	}, nil
}

func (c *Collection) FindByID(
	ctx context.Context,
	id any,
	out any,
) error {
	if err := c.validateOut(out); err != nil {
		return err
	}

	filter, err := idFilter(id)
	if err != nil {
		return err
	}

	err = decodeOneDocument(
		c.native.FindOne(ctx, filter),
		out,
	)
	if err != nil {
		return c.adapter.Normalize(err)
	}

	return nil
}

func (c *Collection) Find(
	ctx context.Context,
	filter docdriver.Filter,
	out any,
	options docdriver.FindOptions,
) error {
	if err := c.validateOut(out); err != nil {
		return err
	}

	mongoFilter, err := buildMongoFilter(filter)
	if err != nil {
		return err
	}

	findOptions := options2FindOptions(options)

	cursor, err := c.native.Find(
		ctx,
		mongoFilter,
		findOptions,
	)
	if err != nil {
		return c.adapter.Normalize(err)
	}
	defer cursor.Close(ctx)

	if err := decodeManyDocuments(ctx, cursor, out); err != nil {
		return c.adapter.Normalize(err)
	}

	return nil
}

func (c *Collection) UpdateByID(
	ctx context.Context,
	id any,
	rev string,
	document any,
) (docdriver.UpdateResult, error) {
	if err := c.validate(document); err != nil {
		return docdriver.UpdateResult{}, err
	}

	filter, err := idFilter(id)
	if err != nil {
		return docdriver.UpdateResult{}, err
	}

	result, err := c.native.ReplaceOne(
		ctx,
		filter,
		document,
	)
	if err != nil {
		return docdriver.UpdateResult{},
			c.adapter.Normalize(err)
	}

	if result.MatchedCount == 0 {
		return docdriver.UpdateResult{},
			c.adapter.Normalize(mongodriver.ErrNoDocuments)
	}

	setDocumentID(document, id)

	return docdriver.UpdateResult{
		Matched:  result.MatchedCount,
		Modified: result.ModifiedCount,
	}, nil
}

func (c *Collection) DeleteByID(
	ctx context.Context,
	id any,
	rev string,
) (docdriver.DeleteResult, error) {
	if c == nil || c.native == nil {
		return docdriver.DeleteResult{}, fmt.Errorf(
			"gorix mongo: collection is unavailable",
		)
	}

	filter, err := idFilter(id)
	if err != nil {
		return docdriver.DeleteResult{}, err
	}

	result, err := c.native.DeleteOne(ctx, filter)
	if err != nil {
		return docdriver.DeleteResult{},
			c.adapter.Normalize(err)
	}

	if result.DeletedCount == 0 {
		return docdriver.DeleteResult{},
			c.adapter.Normalize(mongodriver.ErrNoDocuments)
	}

	return docdriver.DeleteResult{
		Deleted: result.DeletedCount,
	}, nil
}

func (c *Collection) Count(
	ctx context.Context,
	filter docdriver.Filter,
) (int64, error) {
	if c == nil || c.native == nil {
		return 0, fmt.Errorf(
			"gorix mongo: collection is unavailable",
		)
	}

	mongoFilter, err := buildMongoFilter(filter)
	if err != nil {
		return 0, err
	}

	count, err := c.native.CountDocuments(
		ctx,
		mongoFilter,
	)
	if err != nil {
		return 0, c.adapter.Normalize(err)
	}

	return count, nil
}

func (c *Collection) validate(document any) error {
	if c == nil || c.native == nil {
		return fmt.Errorf(
			"gorix mongo: collection is unavailable",
		)
	}

	if document == nil {
		return fmt.Errorf(
			"gorix mongo: document cannot be nil",
		)
	}

	return nil
}

func (c *Collection) validateOut(out any) error {
	if c == nil || c.native == nil {
		return fmt.Errorf(
			"gorix mongo: collection is unavailable",
		)
	}

	if out == nil {
		return fmt.Errorf(
			"gorix mongo: output cannot be nil",
		)
	}

	value := reflect.ValueOf(out)
	if value.Kind() != reflect.Pointer || value.IsNil() {
		return fmt.Errorf(
			"gorix mongo: output must be a non-nil pointer",
		)
	}

	return nil
}

func idFilter(id any) (bson.M, error) {
	if id == nil {
		return nil, fmt.Errorf(
			"gorix mongo: document id cannot be nil",
		)
	}

	if text, ok := id.(string); ok {
		text = strings.TrimSpace(text)

		if text == "" {
			return nil, fmt.Errorf(
				"gorix mongo: document id cannot be empty",
			)
		}

		objectID, err := bson.ObjectIDFromHex(text)
		if err == nil {
			return bson.M{
				"_id": bson.M{
					"$in": []any{
						objectID,
						text,
					},
				},
			}, nil
		}

		return bson.M{"_id": text}, nil
	}

	return bson.M{"_id": id}, nil
}

func options2FindOptions(
	findOptions docdriver.FindOptions,
) *options.FindOptionsBuilder {
	builder := options.Find()

	if findOptions.Limit > 0 {
		builder.SetLimit(findOptions.Limit)
	}

	if findOptions.Offset > 0 {
		builder.SetSkip(findOptions.Offset)
	}

	if len(findOptions.Sort) > 0 {
		sort := bson.D{}

		for _, field := range findOptions.Sort {
			direction := int32(1)
			if field.Desc {
				direction = -1
			}

			sort = append(
				sort,
				bson.E{
					Key:   field.Field,
					Value: direction,
				},
			)
		}

		builder.SetSort(sort)
	}

	return builder
}
