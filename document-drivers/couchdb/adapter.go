package couchdb

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	docdriver "github.com/Gromosome/gorix/document-driver-manager"
	kivik "github.com/go-kivik/kivik/v4"
	_ "github.com/go-kivik/kivik/v4/couchdb"
)

const DriverName = "couchdb"

type Adapter struct{}

func init() {
	docdriver.Register(Adapter{})
}

func (Adapter) Name() string {
	return DriverName
}

func (a Adapter) Open(
	ctx context.Context,
	config docdriver.Config,
) (docdriver.Client, error) {
	if strings.TrimSpace(config.DSN) == "" {
		return nil, fmt.Errorf("gorix couchdb: DSN is required")
	}

	native, err := kivik.New("couch", config.DSN)
	if err != nil {
		return nil, a.Normalize(err)
	}

	client := &Client{
		native:  native,
		adapter: a,
	}

	if err := client.Ping(ctx); err != nil {
		_ = native.Close()
		return nil, err
	}

	return client, nil
}

func (Adapter) Normalize(err error) *docdriver.Error {
	if err == nil {
		return nil
	}

	if dbErr, ok := docdriver.AsError(err); ok {
		return dbErr
	}

	if errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded) {
		return &docdriver.Error{
			Kind:    docdriver.ErrorTimeout,
			Driver:  DriverName,
			Message: err.Error(),
			Cause:   err,
		}
	}

	status := kivik.HTTPStatus(err)
	code := ""
	if status > 0 {
		code = strconv.Itoa(status)
	}

	kind := docdriver.ErrorUnknown

	switch status {
	case http.StatusNotFound:
		kind = docdriver.ErrorNotFound

	case http.StatusConflict:
		kind = docdriver.ErrorConflict

	case http.StatusBadRequest, http.StatusUnprocessableEntity:
		kind = docdriver.ErrorValidation

	case http.StatusUnauthorized, http.StatusForbidden:
		kind = docdriver.ErrorPermissionDenied

	case http.StatusRequestTimeout, http.StatusGatewayTimeout:
		kind = docdriver.ErrorTimeout

	case http.StatusBadGateway, http.StatusServiceUnavailable:
		kind = docdriver.ErrorConnection
	}

	return &docdriver.Error{
		Kind:    kind,
		Driver:  DriverName,
		Code:    code,
		Message: err.Error(),
		Cause:   err,
	}
}
