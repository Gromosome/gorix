package couchdb

import (
	"context"
	"fmt"

	docdriver "github.com/Gromosome/gorix/document-driver-manager"
	kivik "github.com/go-kivik/kivik/v4"
)

type Client struct {
	native  *kivik.Client
	adapter docdriver.Adapter
}

func (c *Client) Ping(ctx context.Context) error {
	if c == nil || c.native == nil {
		return fmt.Errorf("gorix couchdb: client is unavailable")
	}

	ok, err := c.native.Ping(ctx)
	if err != nil {
		return c.adapter.Normalize(err)
	}

	if !ok {
		return fmt.Errorf("gorix couchdb: server is unavailable")
	}

	return nil
}

func (c *Client) Close(ctx context.Context) error {
	if c == nil || c.native == nil {
		return nil
	}

	err := c.native.Close()
	if err != nil {
		return c.adapter.Normalize(err)
	}

	return nil
}

func (c *Client) Database(name string) docdriver.Database {
	return &Database{
		client: c,
		name:   normalizeDatabaseName(name),
	}
}
