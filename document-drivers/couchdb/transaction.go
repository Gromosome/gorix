package couchdb

import (
	"context"

	docdriver "github.com/Gromosome/gorix/document-driver-manager"
)

func (d *Database) BeginTx(
	ctx context.Context,
	options docdriver.TxOptions,
) (docdriver.Tx, error) {
	return nil, &docdriver.Error{
		Kind:     docdriver.ErrorTransactionUnsupported,
		Driver:   DriverName,
		Database: d.Name(),
		Message:  "gorix couchdb: multi-document transactions are not supported; use document revision conflict control instead",
	}
}
