package mysql

import (
	"errors"
	"strconv"

	sqldriver "github.com/Gromosome/gorix/sql-driver-manager"
	mysql "github.com/go-sql-driver/mysql"
)

type adapter struct{}

func (adapter) Name() string          { return "mysql" }
func (adapter) SQLDriverName() string { return "mysql" }

func (adapter) Normalize(err error) *sqldriver.Error {
	if err == nil {
		return nil
	}

	var target *mysql.MySQLError
	if !errors.As(err, &target) {
		return sqldriver.GenericError("mysql", err)
	}

	kinds := map[uint16]sqldriver.ErrorKind{
		1062: sqldriver.ErrorDuplicateKey,
		1451: sqldriver.ErrorForeignKey,
		1452: sqldriver.ErrorForeignKey,
		1048: sqldriver.ErrorNotNull,
		3819: sqldriver.ErrorCheckConstraint,
		1213: sqldriver.ErrorDeadlock,
		1205: sqldriver.ErrorTimeout,
		1064: sqldriver.ErrorSyntax,
		1044: sqldriver.ErrorPermission,
		1045: sqldriver.ErrorPermission,
		1146: sqldriver.ErrorTableNotFound,
		1054: sqldriver.ErrorColumnNotFound,
		2002: sqldriver.ErrorConnection,
		2006: sqldriver.ErrorConnection,
		2013: sqldriver.ErrorConnection,
	}

	kind := kinds[target.Number]
	if kind == "" {
		kind = sqldriver.ErrorUnknown
	}

	return &sqldriver.Error{
		Kind:    kind,
		Driver:  "mysql",
		Code:    strconv.Itoa(int(target.Number)),
		Message: target.Message,
		Cause:   err,
	}
}

func init() {
	sqldriver.Register(adapter{})
}
