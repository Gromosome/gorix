package sqlite_modern

import (
	"errors"
	"strconv"

	sqldriver "github.com/Gromosome/gorix/sql-driver-manager"
	sqlite "modernc.org/sqlite"
)

type adapter struct{}

func (adapter) Name() string          { return "sqlite-modern" }
func (adapter) SQLDriverName() string { return "sqlite" }

func (adapter) Normalize(err error) *sqldriver.Error {
	if err == nil {
		return nil
	}

	var target *sqlite.Error
	if !errors.As(err, &target) {
		return sqldriver.GenericError("sqlite-modern", err)
	}

	code := target.Code()
	baseCode := code & 0xff
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

	kind := kinds[code]
	if kind == "" {
		kind = kinds[baseCode]
	}
	if kind == "" {
		kind = sqldriver.ErrorUnknown
	}

	return &sqldriver.Error{
		Kind:    kind,
		Driver:  "sqlite-modern",
		Code:    strconv.Itoa(code),
		Message: target.Error(),
		Cause:   err,
	}
}

func init() {
	sqldriver.Register(adapter{})
}
