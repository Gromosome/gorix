package document_driver_manager

import (
	"fmt"
	"strings"
)

type Operator string

const (
	OperatorEq     Operator = "eq"
	OperatorNotEq  Operator = "ne"
	OperatorGt     Operator = "gt"
	OperatorGte    Operator = "gte"
	OperatorLt     Operator = "lt"
	OperatorLte    Operator = "lte"
	OperatorIn     Operator = "in"
	OperatorNotIn  Operator = "nin"
	OperatorExists Operator = "exists"
)

type FieldFilter map[Operator]any

type Condition struct {
	Field    string
	Operator Operator
	Value    any
}

type Query struct {
	conditions []Condition
	options    FindOptions
	errors     []error
}

func NewQuery() *Query {
	return &Query{
		conditions: make([]Condition, 0),
		options:    FindOptions{},
		errors:     make([]error, 0),
	}
}

func (q *Query) Where(
	field string,
	operator Operator,
	value any,
) *Query {
	if q == nil {
		return q
	}

	field = strings.TrimSpace(field)

	if field == "" {
		q.errors = append(
			q.errors,
			fmt.Errorf("gorix document query: field cannot be empty"),
		)
		return q
	}

	if operator == "" {
		q.errors = append(
			q.errors,
			fmt.Errorf(
				"gorix document query: operator cannot be empty for field %q",
				field,
			),
		)
		return q
	}

	q.conditions = append(
		q.conditions,
		Condition{
			Field:    field,
			Operator: operator,
			Value:    value,
		},
	)

	return q
}

func (q *Query) Eq(field string, value any) *Query {
	return q.Where(field, OperatorEq, value)
}

func (q *Query) NotEq(field string, value any) *Query {
	return q.Where(field, OperatorNotEq, value)
}

func (q *Query) Gt(field string, value any) *Query {
	return q.Where(field, OperatorGt, value)
}

func (q *Query) Gte(field string, value any) *Query {
	return q.Where(field, OperatorGte, value)
}

func (q *Query) Lt(field string, value any) *Query {
	return q.Where(field, OperatorLt, value)
}

func (q *Query) Lte(field string, value any) *Query {
	return q.Where(field, OperatorLte, value)
}

func (q *Query) In(field string, values ...any) *Query {
	return q.Where(field, OperatorIn, values)
}

func (q *Query) NotIn(field string, values ...any) *Query {
	return q.Where(field, OperatorNotIn, values)
}

func (q *Query) Exists(field string, exists bool) *Query {
	return q.Where(field, OperatorExists, exists)
}

func (q *Query) Limit(limit int64) *Query {
	if q == nil {
		return q
	}

	if limit < 0 {
		q.errors = append(
			q.errors,
			fmt.Errorf("gorix document query: limit cannot be negative"),
		)
		return q
	}

	q.options.Limit = limit
	return q
}

func (q *Query) Offset(offset int64) *Query {
	if q == nil {
		return q
	}

	if offset < 0 {
		q.errors = append(
			q.errors,
			fmt.Errorf("gorix document query: offset cannot be negative"),
		)
		return q
	}

	q.options.Offset = offset
	return q
}

func (q *Query) SortAsc(field string) *Query {
	return q.sort(field, false)
}

func (q *Query) SortDesc(field string) *Query {
	return q.sort(field, true)
}

func (q *Query) sort(field string, desc bool) *Query {
	if q == nil {
		return q
	}

	field = strings.TrimSpace(field)

	if field == "" {
		q.errors = append(
			q.errors,
			fmt.Errorf("gorix document query: sort field cannot be empty"),
		)
		return q
	}

	q.options.Sort = append(
		q.options.Sort,
		SortField{
			Field: field,
			Desc:  desc,
		},
	)

	return q
}

func (q *Query) Build() (
	Filter,
	FindOptions,
	error,
) {
	if q == nil {
		return nil, FindOptions{}, fmt.Errorf(
			"gorix document query: query cannot be nil",
		)
	}

	if len(q.errors) > 0 {
		return nil, FindOptions{}, q.errors[0]
	}

	filter := make(Filter)

	for _, condition := range q.conditions {
		current, exists := filter[condition.Field]

		if !exists {
			filter[condition.Field] = FieldFilter{
				condition.Operator: condition.Value,
			}
			continue
		}

		fieldFilter, ok := current.(FieldFilter)
		if !ok {
			return nil, FindOptions{}, fmt.Errorf(
				"gorix document query: invalid filter state for field %q",
				condition.Field,
			)
		}

		fieldFilter[condition.Operator] = condition.Value
		filter[condition.Field] = fieldFilter
	}

	return filter, q.options, nil
}
