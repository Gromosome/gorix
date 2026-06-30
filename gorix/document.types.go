package gorix

import (
	docdriver "github.com/Gromosome/gorix/document-driver-manager"
	"github.com/Gromosome/gorix/gorix/config"
	"github.com/Gromosome/gorix/gorix/core/document"
)

type DocumentManager = document.Manager
type DocumentConfig = config.Config

type DocumentRepository[T any, ID comparable] = document.Repository[T, ID]
type DocumentScopedRepository[T any, ID comparable] = document.ScopedRepository[T, ID]

type DocumentMetadata = document.DocumentMetadata
type DocumentSchema = document.DocumentSchema
type DocumentIndex = document.DocumentIndex

type DocumentTx = document.Tx
type DocumentTxOptions = document.TxOptions
type DocumentTransactionFunc = document.TransactionFunc

type DocumentFilter = docdriver.Filter
type DocumentFindOptions = docdriver.FindOptions
type DocumentQuery = docdriver.Query

var NewDocumentQuery = docdriver.NewQuery
var WithDocumentTransaction = document.WithTransaction
var AsDocumentError = document.AsError
var IsDocumentErrorKind = document.IsKind

func NewDocumentRepository[T any, ID comparable](
	manager *document.Manager,
	connectionNames ...string,
) (*DocumentRepository[T, ID], error) {
	return document.NewRepository[T, ID](manager, connectionNames...)
}
