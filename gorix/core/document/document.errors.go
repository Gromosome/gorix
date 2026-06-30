package document

import (
	"fmt"

	docdriver "github.com/Gromosome/gorix/document-driver-manager"
	"github.com/Gromosome/gorix/gorix/config"
)

const (
	ErrorInvalidMetadata docdriver.ErrorKind = "invalid_metadata"
)

func NewError(
	kind docdriver.ErrorKind,
	message string,
	cause error,
) *docdriver.Error {
	return &docdriver.Error{
		Kind:    kind,
		Driver:  "gorix-document",
		Message: message,
		Cause:   cause,
	}
}

func AsError(err error) (*docdriver.Error, bool) {
	return docdriver.AsError(err)
}

func IsKind(err error, kind docdriver.ErrorKind) bool {
	return docdriver.IsKind(err, kind)
}

func TransactionUnsupportedError(connectionName string) *docdriver.Error {
	if connectionName == "" {
		connectionName = config.DefaultConnectionName
	}

	return &docdriver.Error{
		Kind:     docdriver.ErrorTransactionUnsupported,
		Driver:   "gorix-document",
		Database: connectionName,
		Message: fmt.Sprintf(
			"gorix document: connection %q does not support document transactions",
			connectionName,
		),
	}
}
