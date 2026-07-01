package mongo

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	docdriver "github.com/Gromosome/gorix/document-driver-manager"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

const DriverName = "mongo"

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
		return nil, fmt.Errorf(
			"gorix mongo: DSN is required",
		)
	}

	native, err := mongodriver.Connect(
		options.Client().ApplyURI(config.DSN),
	)
	if err != nil {
		return nil, err
	}

	if err := native.Ping(ctx, readpref.Primary()); err != nil {
		_ = native.Disconnect(ctx)
		return nil, err
	}

	return &Client{
		native:  native,
		adapter: a,
	}, nil
}

func (Adapter) Normalize(err error) *docdriver.Error {
	if err == nil {
		return nil
	}

	if dbErr, ok := docdriver.AsError(err); ok {
		return dbErr
	}

	if errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded) ||
		mongodriver.IsTimeout(err) {
		return &docdriver.Error{
			Kind:    docdriver.ErrorTimeout,
			Driver:  DriverName,
			Message: err.Error(),
			Cause:   err,
		}
	}

	if errors.Is(err, mongodriver.ErrNoDocuments) {
		return &docdriver.Error{
			Kind:    docdriver.ErrorNotFound,
			Driver:  DriverName,
			Message: err.Error(),
			Cause:   err,
		}
	}

	if mongodriver.IsDuplicateKeyError(err) {
		return &docdriver.Error{
			Kind:    docdriver.ErrorDuplicateKey,
			Driver:  DriverName,
			Code:    "11000",
			Message: err.Error(),
			Cause:   err,
		}
	}

	if mongodriver.IsNetworkError(err) {
		return &docdriver.Error{
			Kind:    docdriver.ErrorConnection,
			Driver:  DriverName,
			Message: err.Error(),
			Cause:   err,
		}
	}

	code := firstErrorCode(err)

	return &docdriver.Error{
		Kind:    docdriver.ErrorUnknown,
		Driver:  DriverName,
		Code:    code,
		Message: err.Error(),
		Cause:   err,
	}
}

func firstErrorCode(err error) string {
	codes := mongodriver.ErrorCodes(err)
	if len(codes) == 0 {
		return ""
	}

	return strconv.Itoa(codes[0])
}
