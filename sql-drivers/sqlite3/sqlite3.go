package sqlite3

import (
	"errors"
	"strconv"

	sqldriver "github.com/Gromosome/gorix/sql-driver-manager"
	sqlite3 "github.com/mattn/go-sqlite3"
)

type adapter struct{}

func (adapter) Name() string          { return "sqlite3" }
func (adapter) SQLDriverName() string { return "sqlite3" }

func (adapter) Normalize(err error) *sqldriver.Error {
	if err == nil {
		return nil
	}

	var target sqlite3.Error
	if !errors.As(err, &target) {
		return sqldriver.GenericError("sqlite3", err)
	}

	extendedCode := int(target.ExtendedCode)
	kinds := map[int]sqldriver.ErrorKind{
		1555: sqldriver.ErrorDuplicateKey,
		2067: sqldriver.ErrorUniqueViolation,
		787:  sqldriver.ErrorForeignKey,
		1299: sqldriver.ErrorNotNull,
		275:  sqldriver.ErrorCheckConstraint,
		5:    sqldriver.ErrorTimeout,
		6:    sqldriver.ErrorTimeout,
		1:    sqldriver.ErrorSyntax,
		23:   sqldriver.ErrorPermission,
	}

	kind := kinds[extendedCode]
	if kind == "" {
		kind = kinds[int(target.Code)]
	}
	if kind == "" {
		kind = sqldriver.ErrorUnknown
	}

	return &sqldriver.Error{
		Kind:    kind,
		Driver:  "sqlite3",
		Code:    strconv.Itoa(extendedCode),
		Message: target.Error(),
		Cause:   err,
	}
}

func init() {
	sqldriver.Register(adapter{})
}
