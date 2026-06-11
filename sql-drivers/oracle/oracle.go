package oracle

import (
	"errors"
	"strconv"

	sqldriver "github.com/Gromosome/gorix/sql-driver-manager"
	_ "github.com/godror/godror"
)

type oracleError interface {
	error
	Code() int
}

type adapter struct{}

func (adapter) Name() string          { return "oracle" }
func (adapter) SQLDriverName() string { return "godror" }

func (adapter) Normalize(err error) *sqldriver.Error {
	if err == nil {
		return nil
	}

	var target oracleError
	if !errors.As(err, &target) {
		return sqldriver.GenericError("oracle", err)
	}

	code := target.Code()
	kinds := map[int]sqldriver.ErrorKind{
		1:     sqldriver.ErrorUniqueViolation,
		2291:  sqldriver.ErrorForeignKey,
		2292:  sqldriver.ErrorForeignKey,
		1400:  sqldriver.ErrorNotNull,
		2290:  sqldriver.ErrorCheckConstraint,
		8177:  sqldriver.ErrorSerialization,
		60:    sqldriver.ErrorDeadlock,
		1013:  sqldriver.ErrorTimeout,
		900:   sqldriver.ErrorSyntax,
		903:   sqldriver.ErrorSyntax,
		942:   sqldriver.ErrorTableNotFound,
		904:   sqldriver.ErrorColumnNotFound,
		1031:  sqldriver.ErrorPermission,
		1017:  sqldriver.ErrorPermission,
		3113:  sqldriver.ErrorConnection,
		3114:  sqldriver.ErrorConnection,
		12170: sqldriver.ErrorTimeout,
	}

	kind := kinds[code]
	if kind == "" {
		kind = sqldriver.ErrorUnknown
	}

	return &sqldriver.Error{
		Kind:    kind,
		Driver:  "oracle",
		Code:    strconv.Itoa(code),
		Message: target.Error(),
		Cause:   err,
	}
}

func init() {
	sqldriver.Register(adapter{})
}
