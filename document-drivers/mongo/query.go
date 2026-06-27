package mongo

import (
	"fmt"

	docdriver "github.com/Gromosome/gorix/document-driver-manager"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func buildMongoFilter(
	filter docdriver.Filter,
) (bson.M, error) {
	mongoFilter := bson.M{}

	for field, value := range filter {
		converted, err := convertFilterValue(value)
		if err != nil {
			return nil, fmt.Errorf(
				"gorix mongo: invalid filter for field %q: %w",
				field,
				err,
			)
		}

		mongoField := field
		if field == "id" {
			mongoField = "_id"
		}

		mongoFilter[mongoField] = converted
	}

	return mongoFilter, nil
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

	selector := bson.M{}

	for operator, operatorValue := range fieldFilter {
		mongoOperator, err := toMongoOperator(operator)
		if err != nil {
			return nil, err
		}

		selector[mongoOperator] = operatorValue
	}

	return selector, nil
}

func toMongoOperator(
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
