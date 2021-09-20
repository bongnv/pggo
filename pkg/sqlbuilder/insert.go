package sqlbuilder

import (
	"errors"
	"strings"
)

// Insert starts a new INSERT query.
func Insert(table Table) *InsertBuilder {
	return &InsertBuilder{
		table: table,
	}
}

// Insert starts a new INSERT query with a table name.
func InsertTable(tableName string) *InsertBuilder {
	return &InsertBuilder{
		table: BaseTable(tableName),
	}
}

// InsertBuilder is a builder to build an INSERT query.
type InsertBuilder struct {
	cols   []string
	table  Table
	values [][]interface{}
}

// Columns adds columns to the INSERT query.
func (b *InsertBuilder) Columns(cols ...string) *InsertBuilder {
	b.cols = cols
	return b
}

// Values adds a single row's values to the query.
// Multiple calls will create multiple rows to the query.
func (b *InsertBuilder) Values(values ...interface{}) *InsertBuilder {
	b.values = append(b.values, values)
	return b
}

// SQL compiles all provided data to return an INSERT query and arguments.
func (b InsertBuilder) SQL() (string, []interface{}, error) {
	if len(b.values) == 0 {
		return "", nil, errors.New("sqlbuilder: there must be at least one row")
	}

	sb := &strings.Builder{}
	args := argumentList{}

	_, _ = sb.WriteString("INSERT INTO ")

	b.table.Build(sb, &args)

	_, _ = sb.WriteString(" (")
	for i, col := range b.cols {
		if i > 0 {
			_, _ = sb.WriteString(",")
		}
		_, _ = sb.WriteString(col)
	}
	_, _ = sb.WriteString(") VALUES ")

	for i, values := range b.values {
		if i > 0 {
			_, _ = sb.WriteString(",")
		}
		groupPlaceholder{
			values: values,
		}.Build(sb, &args)
	}

	return sb.String(), args, nil
}
