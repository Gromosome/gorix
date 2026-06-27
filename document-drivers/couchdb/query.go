package couchdb

import (
	"fmt"

	docdriver "github.com/Gromosome/gorix/document-driver-manager"
)

type mangoQuery struct {
	Selector map[string]any      `json:"selector"`
	Limit    int64               `json:"limit,omitempty"`
	Skip     int64               `json:"skip,omitempty"`
	Sort     []map[string]string `json:"sort,omitempty"`
	Fields   []string            `json:"fields,omitempty"`
}

func buildMangoQuery(
	filter docdriver.Filter,
	options docdriver.FindOptions,
) (mangoQuery, error) {
	query := mangoQuery{
		Selector: make(map[string]any),
	}

	for field, value := range filter {
		converted, err := convertFilterValue(value)
		if err != nil {
			return mangoQuery{}, fmt.Errorf(
				"gorix couchdb: invalid filter for field %q: %w",
				field,
				err,
			)
		}

		query.Selector[field] = converted
	}

	if len(query.Selector) == 0 {
		query.Selector["_id"] = map[string]any{
			"$exists": true,
		}
	}

	if options.Limit > 0 {
		query.Limit = options.Limit
	}

	if options.Offset > 0 {
		query.Skip = options.Offset
	}

	for _, sort := range options.Sort {
		direction := "asc"
		if sort.Desc {
			direction = "desc"
		}

		query.Sort = append(
			query.Sort,
			map[string]string{
				sort.Field: direction,
			},
		)
	}

	return query, nil
}

func convertFilterValue(value any) (any, error) {
	fieldFilter, ok := value.(docdriver.FieldFilter)
	if !ok {
		return value, nil
	}

	if len(fieldFilter) == 1 {
		if eqValue, exists := fieldFilter[docdriver.OperatorEq]; exists {
			return eqValue, nil
		}
	}

	selector := make(map[string]any)

	for operator, operatorValue := range fieldFilter {
		mangoOperator, err := toMangoOperator(operator)
		if err != nil {
			return nil, err
		}

		selector[mangoOperator] = operatorValue
	}

	return selector, nil
}

func toMangoOperator(
	operator docdriver.Operator,
) (string, error) {
	switch operator {
	case docdriver.OperatorEq:
		return "$eq", nil

	case docdriver.OperatorNotEq:
		return "$ne", nil

	case docdriver.OperatorGt:
		return "$gt", nil

	case docdriver.OperatorGte:
		return "$gte", nil

	case docdriver.OperatorLt:
		return "$lt", nil

	case docdriver.OperatorLte:
		return "$lte", nil

	case docdriver.OperatorIn:
		return "$in", nil

	case docdriver.OperatorNotIn:
		return "$nin", nil

	case docdriver.OperatorExists:
		return "$exists", nil

	default:
		return "", fmt.Errorf(
			"unsupported operator %q",
			operator,
		)
	}
}
