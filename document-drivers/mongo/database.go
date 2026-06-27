package mongo

import (
	docdriver "github.com/Gromosome/gorix/document-driver-manager"
)

type Database struct {
	client *Client
	name   string
}

func (d *Database) Name() string {
	if d == nil {
		return ""
	}

	return d.name
}

func (d *Database) Collection(name string) docdriver.Collection {
	collectionName := normalizeCollectionName(name)

	return &Collection{
		name: collectionName,
		native: d.client.native.
			Database(d.name).
			Collection(collectionName),
		adapter: d.client.adapter,
	}
}
