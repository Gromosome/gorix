package document_driver_manager

type Filter map[string]any

type FindOptions struct {
	Limit  int64
	Offset int64
	Sort   []SortField
}

type SortField struct {
	Field string
	Desc  bool
}

type InsertResult struct {
	ID  any
	Rev string
}

type UpdateResult struct {
	Matched  int64
	Modified int64
	Rev      string
}

type DeleteResult struct {
	Deleted int64
}
