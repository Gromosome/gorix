package mongo

import (
	"context"
	"fmt"
	"sync"

	docdriver "github.com/Gromosome/gorix/document-driver-manager"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

type Tx struct {
	database *Database
	session  *mongodriver.Session
	mutex    sync.Mutex
	closed   bool
}

func (d *Database) BeginTx(
	ctx context.Context,
	options docdriver.TxOptions,
) (docdriver.Tx, error) {
	_ = options

	if d == nil || d.client == nil || d.client.native == nil {
		return nil, fmt.Errorf("gorix mongo: database is unavailable")
	}

	session, err := d.client.native.StartSession()
	if err != nil {
		return nil, d.client.adapter.Normalize(err)
	}

	if err := session.StartTransaction(); err != nil {
		session.EndSession(ctx)
		return nil, d.client.adapter.Normalize(err)
	}

	return &Tx{
		database: d,
		session:  session,
	}, nil
}

func (tx *Tx) Collection(name string) docdriver.Collection {
	collectionName := normalizeCollectionName(name)

	base := &Collection{
		name: collectionName,
		native: tx.database.client.native.
			Database(tx.database.name).
			Collection(collectionName),
		adapter: tx.database.client.adapter,
	}

	return &TxCollection{
		base:    base,
		session: tx.session,
	}
}

func (tx *Tx) Commit(ctx context.Context) error {
	if tx == nil || tx.session == nil {
		return nil
	}

	tx.mutex.Lock()
	defer tx.mutex.Unlock()

	if tx.closed {
		return nil
	}

	tx.closed = true
	defer tx.session.EndSession(ctx)

	err := tx.session.CommitTransaction(ctx)
	if err != nil {
		return tx.database.client.adapter.Normalize(err)
	}

	return nil
}

func (tx *Tx) Rollback(ctx context.Context) error {
	if tx == nil || tx.session == nil {
		return nil
	}

	tx.mutex.Lock()
	defer tx.mutex.Unlock()

	if tx.closed {
		return nil
	}

	tx.closed = true
	defer tx.session.EndSession(ctx)

	err := tx.session.AbortTransaction(ctx)
	if err != nil {
		return tx.database.client.adapter.Normalize(err)
	}

	return nil
}

type TxCollection struct {
	base    *Collection
	session *mongodriver.Session
}

func (c *TxCollection) Name() string {
	return c.base.Name()
}

func (c *TxCollection) InsertOne(
	ctx context.Context,
	document any,
) (docdriver.InsertResult, error) {
	return c.base.InsertOne(
		mongodriver.NewSessionContext(ctx, c.session),
		document,
	)
}

func (c *TxCollection) FindByID(
	ctx context.Context,
	id any,
	out any,
) error {
	return c.base.FindByID(
		mongodriver.NewSessionContext(ctx, c.session),
		id,
		out,
	)
}

func (c *TxCollection) Find(
	ctx context.Context,
	filter docdriver.Filter,
	out any,
	options docdriver.FindOptions,
) error {
	return c.base.Find(
		mongodriver.NewSessionContext(ctx, c.session),
		filter,
		out,
		options,
	)
}

func (c *TxCollection) UpdateByID(
	ctx context.Context,
	id any,
	rev string,
	document any,
) (docdriver.UpdateResult, error) {
	return c.base.UpdateByID(
		mongodriver.NewSessionContext(ctx, c.session),
		id,
		rev,
		document,
	)
}

func (c *TxCollection) DeleteByID(
	ctx context.Context,
	id any,
	rev string,
) (docdriver.DeleteResult, error) {
	return c.base.DeleteByID(
		mongodriver.NewSessionContext(ctx, c.session),
		id,
		rev,
	)
}

func (c *TxCollection) Count(
	ctx context.Context,
	filter docdriver.Filter,
) (int64, error) {
	return c.base.Count(
		mongodriver.NewSessionContext(ctx, c.session),
		filter,
	)
}
