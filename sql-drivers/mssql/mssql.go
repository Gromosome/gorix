package mssql

import (
	"errors"
	"strconv"

	sqldriver "github.com/Gromosome/gorix/sql-driver-manager"
	mssql "github.com/microsoft/go-mssqldb"
)

type adapter struct{}

func (adapter) Name() string          { return "mssql" }
func (adapter) SQLDriverName() string { return "sqlserver" }

func (adapter) Normalize(err error) *sqldriver.Error {
	if err == nil {
		return nil
	}

	var target mssql.Error
	if !errors.As(err, &target) {
		return sqldriver.GenericError("mssql", err)
	}

	kinds := map[int32]sqldriver.ErrorKind{
		2601:  sqldriver.ErrorUniqueViolation,
		2627:  sqldriver.ErrorDuplicateKey,
		547:   sqldriver.ErrorForeignKey,
		515:   sqldriver.ErrorNotNull,
		1205:  sqldriver.ErrorDeadlock,
		1222:  sqldriver.ErrorTimeout,
		102:   sqldriver.ErrorSyntax,
		156:   sqldriver.ErrorSyntax,
		208:   sqldriver.ErrorTableNotFound,
		207:   sqldriver.ErrorColumnNotFound,
		229:   sqldriver.ErrorPermission,
		18456: sqldriver.ErrorPermission,
		53:    sqldriver.ErrorConnection,
		233:   sqldriver.ErrorConnection,
	}

	kind := kinds[target.Number]
	if kind == "" {
		kind = sqldriver.ErrorUnknown
	}

	return &sqldriver.Error{
		Kind:    kind,
		Driver:  "mssql",
		Code:    strconv.FormatInt(int64(target.Number), 10),
		Message: target.Message,
		Cause:   err,
	}
}

func init() {
	sqldriver.Register(adapter{})
}
