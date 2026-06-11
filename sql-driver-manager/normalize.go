package sql_driver_manager

import (
	"context"
	"errors"
	"net"
	"os"
	"strings"
)

func GenericError(driver string, err error) *Error {
	if err == nil {
		return nil
	}

	kind := ErrorUnknown

	if errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, os.ErrDeadlineExceeded) {
		kind = ErrorTimeout
	}

	var networkError net.Error
	if errors.As(err, &networkError) {
		if networkError.Timeout() {
			kind = ErrorTimeout
		} else {
			kind = ErrorConnection
		}
	}

	message := strings.ToLower(err.Error())
	if kind == ErrorUnknown &&
		(strings.Contains(message, "connection refused") ||
			strings.Contains(message, "connection reset") ||
			strings.Contains(message, "broken pipe") ||
			strings.Contains(message, "bad connection")) {
		kind = ErrorConnection
	}

	return &Error{
		Kind:    kind,
		Driver:  driver,
		Message: err.Error(),
		Cause:   err,
	}
}
