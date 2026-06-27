package couchdb

import (
	"strings"

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
	collectionName := normalizeDatabaseName(name)

	return &Collection{
		name: collectionName,
		db: d.client.native.DB(
			d.collectionDatabaseName(collectionName),
		),
		adapter: d.client.adapter,
	}
}

func (d *Database) collectionDatabaseName(
	collection string,
) string {
	base := normalizeDatabaseName(d.name)

	if strings.TrimSpace(collection) == "" {
		return base
	}

	return normalizeDatabaseName(
		base + "_" + collection,
	)
}
