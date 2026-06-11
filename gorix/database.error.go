package gorix

import sqldriver "github.com/Gromosome/gorix/sql-driver-manager"

type DatabaseError = sqldriver.Error
type DatabaseErrorKind = sqldriver.ErrorKind

const (
	DatabaseErrorUnknown = sqldriver.ErrorUnknown

	DatabaseErrorDuplicateKey = sqldriver.ErrorDuplicateKey

	DatabaseErrorUniqueViolation = sqldriver.ErrorUniqueViolation

	DatabaseErrorForeignKey = sqldriver.ErrorForeignKey

	DatabaseErrorNotNull = sqldriver.ErrorNotNull

	DatabaseErrorCheckConstraint = sqldriver.ErrorCheckConstraint

	DatabaseErrorSerialization = sqldriver.ErrorSerialization

	DatabaseErrorDeadlock = sqldriver.ErrorDeadlock

	DatabaseErrorConnection = sqldriver.ErrorConnection

	DatabaseErrorTimeout = sqldriver.ErrorTimeout

	DatabaseErrorSyntax = sqldriver.ErrorSyntax

	DatabaseErrorPermission = sqldriver.ErrorPermission

	DatabaseErrorTableNotFound = sqldriver.ErrorTableNotFound

	DatabaseErrorColumnNotFound = sqldriver.ErrorColumnNotFound
)

func AsDatabaseError(
	err error,
) (*DatabaseError, bool) {
	return sqldriver.AsError(err)
}

func IsDatabaseErrorKind(
	err error,
	kind DatabaseErrorKind,
) bool {
	return sqldriver.IsKind(err, kind)
}
