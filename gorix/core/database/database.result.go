package database

func (r Result) LastInsertID() (
	int64,
	error,
) {
	return r.native.LastInsertId()
}

func (r Result) RowsAffected() (
	int64,
	error,
) {
	return r.native.RowsAffected()
}
