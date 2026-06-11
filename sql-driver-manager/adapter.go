package sql_driver_manager

type Adapter interface {
	Name() string
	SQLDriverName() string
	Normalize(error) *Error
}
