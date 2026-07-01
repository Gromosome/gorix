package mongo

import (
	"context"
	"fmt"

	docdriver "github.com/Gromosome/gorix/document-driver-manager"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Client struct {
	native  *mongodriver.Client
	adapter Adapter
}

func (c *Client) Ping(ctx context.Context) error {
	if c == nil || c.native == nil {
		return fmt.Errorf(
			"gorix mongo: client is unavailable",
		)
	}

	err := c.native.Ping(ctx, readpref.Primary())
	if err != nil {
		return c.adapter.Normalize(err)
	}

	return nil
}

func (c *Client) Close(ctx context.Context) error {
	if c == nil || c.native == nil {
		return nil
	}

	err := c.native.Disconnect(ctx)
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
