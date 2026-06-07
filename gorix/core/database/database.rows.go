package database

func (r *Rows) Next() bool {
	return r != nil &&
		r.native != nil &&
		r.native.Next()
}

func (r *Rows) Scan(
	destinations ...any,
) error {
	return r.native.Scan(destinations...)
}

func (r *Rows) Columns() (
	[]string,
	error,
) {
	return r.native.Columns()
}

func (r *Rows) Err() error {
	if r == nil || r.native == nil {
		return nil
	}

	return r.native.Err()
}

func (r *Rows) Close() error {
	if r == nil || r.native == nil {
		return nil
	}

	return r.native.Close()
}

func (r *Row) Scan(
	destinations ...any,
) error {
	return r.native.Scan(destinations...)
}
