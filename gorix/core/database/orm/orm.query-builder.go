package orm

import (
	"fmt"
	"strings"
)

type Condition struct {
	Expression string
	Args       []any
}

type QueryBuilder struct {
	dialect    Dialect
	table      string
	columns    []string
	conditions []Condition
	orderBy    []string
	groupBy    []string
	limit      *int
	offset     *int
}

func NewQueryBuilder(
	dialect Dialect,
	table string,
) *QueryBuilder {
	return &QueryBuilder{
		dialect: dialect,
		table:   table,
	}
}

func (b *QueryBuilder) Select(
	columns ...string,
) *QueryBuilder {
	b.columns = append(b.columns, columns...)
	return b
}

func (b *QueryBuilder) Where(
	expression string,
	args ...any,
) *QueryBuilder {
	if strings.TrimSpace(expression) != "" {
		b.conditions = append(
			b.conditions,
			Condition{
				Expression: expression,
				Args:       args,
			},
		)
	}

	return b
}

func (b *QueryBuilder) OrderBy(
	expressions ...string,
) *QueryBuilder {
	b.orderBy = append(b.orderBy, expressions...)
	return b
}

func (b *QueryBuilder) GroupBy(
	columns ...string,
) *QueryBuilder {
	b.groupBy = append(b.groupBy, columns...)
	return b
}

func (b *QueryBuilder) Limit(value int) *QueryBuilder {
	b.limit = &value
	return b
}

func (b *QueryBuilder) Offset(value int) *QueryBuilder {
	b.offset = &value
	return b
}

func (b *QueryBuilder) BuildSelect() (
	string,
	[]any,
	error,
) {
	if b.dialect == nil {
		return "", nil, fmt.Errorf(
			"gorix orm: query dialect cannot be nil",
		)
	}

	if strings.TrimSpace(b.table) == "" {
		return "", nil, fmt.Errorf(
			"gorix orm: table cannot be empty",
		)
	}

	columns := "*"

	if len(b.columns) > 0 {
		quoted := make([]string, 0, len(b.columns))

		for _, column := range b.columns {
			if column == "*" {
				quoted = append(quoted, column)
				continue
			}

			quoted = append(
				quoted,
				b.dialect.QuoteIdentifier(column),
			)
		}

		columns = strings.Join(quoted, ", ")
	}

	var builder strings.Builder

	builder.WriteString("SELECT ")
	builder.WriteString(columns)
	builder.WriteString(" FROM ")
	builder.WriteString(
		b.dialect.QuoteIdentifier(b.table),
	)

	args := make([]any, 0)
	placeholderIndex := 1

	if len(b.conditions) > 0 {
		builder.WriteString(" WHERE ")

		conditionSQL := make([]string, 0, len(b.conditions))

		for _, condition := range b.conditions {
			expression := condition.Expression

			for range condition.Args {
				expression = strings.Replace(
					expression,
					"?",
					b.dialect.Placeholder(placeholderIndex),
					1,
				)

				placeholderIndex++
			}

			conditionSQL = append(
				conditionSQL,
				"("+expression+")",
			)

			args = append(args, condition.Args...)
		}

		builder.WriteString(strings.Join(conditionSQL, " AND "))
	}

	if len(b.groupBy) > 0 {
		builder.WriteString(" GROUP BY ")
		builder.WriteString(strings.Join(b.groupBy, ", "))
	}

	if len(b.orderBy) > 0 {
		builder.WriteString(" ORDER BY ")
		builder.WriteString(strings.Join(b.orderBy, ", "))
	}

	if b.limit != nil {
		builder.WriteString(
			fmt.Sprintf(" LIMIT %d", *b.limit),
		)
	}

	if b.offset != nil {
		builder.WriteString(
			fmt.Sprintf(" OFFSET %d", *b.offset),
		)
	}

	return builder.String(), args, nil
}
