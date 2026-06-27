package couchdb

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	docdriver "github.com/Gromosome/gorix/document-driver-manager"
	kivik "github.com/go-kivik/kivik/v4"
)

type Collection struct {
	name    string
	db      *kivik.DB
	adapter docdriver.Adapter
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

	if id, ok := documentID(document); ok && id != "" {
		rev, err := c.db.Put(ctx, id, document)
		if err != nil {
			return docdriver.InsertResult{},
				c.adapter.Normalize(err)
		}

		setDocumentRevision(document, rev)

		return docdriver.InsertResult{
			ID:  id,
			Rev: rev,
		}, nil
	}

	id, rev, err := c.db.CreateDoc(ctx, document)
	if err != nil {
		return docdriver.InsertResult{},
			c.adapter.Normalize(err)
	}

	setDocumentID(document, id)
	setDocumentRevision(document, rev)

	return docdriver.InsertResult{
		ID:  id,
		Rev: rev,
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

	docID := strings.TrimSpace(
		fmt.Sprint(id),
	)
	if docID == "" {
		return fmt.Errorf(
			"gorix couchdb: document id cannot be empty",
		)
	}

	document := c.db.Get(ctx, docID)
	defer document.Close()

	if err := document.ScanDoc(out); err != nil {
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

	query, err := buildMangoQuery(filter, options)
	if err != nil {
		return err
	}

	rows := c.db.Find(ctx, query)
	defer rows.Close()

	if err := scanResultSet(rows, out); err != nil {
		return c.adapter.Normalize(err)
	}

	if err := rows.Err(); err != nil {
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

	docID := strings.TrimSpace(
		fmt.Sprint(id),
	)
	if docID == "" {
		return docdriver.UpdateResult{}, fmt.Errorf(
			"gorix couchdb: document id cannot be empty",
		)
	}

	if strings.TrimSpace(rev) == "" {
		if entityRev, ok := documentRevision(document); ok {
			rev = entityRev
		}
	}

	if strings.TrimSpace(rev) == "" {
		latestRev, err := c.db.GetRev(ctx, docID)
		if err != nil {
			return docdriver.UpdateResult{},
				c.adapter.Normalize(err)
		}

		rev = latestRev
	}

	newRev, err := c.db.Put(
		ctx,
		docID,
		document,
		kivik.Rev(rev),
	)
	if err != nil {
		return docdriver.UpdateResult{},
			c.adapter.Normalize(err)
	}

	setDocumentID(document, docID)
	setDocumentRevision(document, newRev)

	return docdriver.UpdateResult{
		Matched:  1,
		Modified: 1,
		Rev:      newRev,
	}, nil
}

func (c *Collection) DeleteByID(
	ctx context.Context,
	id any,
	rev string,
) (docdriver.DeleteResult, error) {
	if c == nil || c.db == nil {
		return docdriver.DeleteResult{}, fmt.Errorf(
			"gorix couchdb: collection is unavailable",
		)
	}

	docID := strings.TrimSpace(
		fmt.Sprint(id),
	)
	if docID == "" {
		return docdriver.DeleteResult{}, fmt.Errorf(
			"gorix couchdb: document id cannot be empty",
		)
	}

	if strings.TrimSpace(rev) == "" {
		latestRev, err := c.db.GetRev(ctx, docID)
		if err != nil {
			return docdriver.DeleteResult{},
				c.adapter.Normalize(err)
		}

		rev = latestRev
	}

	_, err := c.db.Delete(ctx, docID, rev)
	if err != nil {
		return docdriver.DeleteResult{},
			c.adapter.Normalize(err)
	}

	return docdriver.DeleteResult{
		Deleted: 1,
	}, nil
}

func (c *Collection) Count(
	ctx context.Context,
	filter docdriver.Filter,
) (int64, error) {
	query, err := buildMangoQuery(
		filter,
		docdriver.FindOptions{},
	)
	if err != nil {
		return 0, err
	}

	query.Fields = []string{"_id"}

	rows := c.db.Find(ctx, query)
	defer rows.Close()

	var count int64

	for rows.Next() {
		var item map[string]any
		if err := rows.ScanDoc(&item); err != nil {
			return 0, c.adapter.Normalize(err)
		}

		count++
	}

	if err := rows.Err(); err != nil {
		return 0, c.adapter.Normalize(err)
	}

	return count, nil
}

func (c *Collection) validate(document any) error {
	if c == nil || c.db == nil {
		return fmt.Errorf(
			"gorix couchdb: collection is unavailable",
		)
	}

	if document == nil {
		return fmt.Errorf(
			"gorix couchdb: document cannot be nil",
		)
	}

	return nil
}

func (c *Collection) validateOut(out any) error {
	if c == nil || c.db == nil {
		return fmt.Errorf(
			"gorix couchdb: collection is unavailable",
		)
	}

	if out == nil {
		return fmt.Errorf(
			"gorix couchdb: output cannot be nil",
		)
	}

	if reflect.ValueOf(out).Kind() != reflect.Pointer {
		return fmt.Errorf(
			"gorix couchdb: output must be a pointer",
		)
	}

	return nil
}
