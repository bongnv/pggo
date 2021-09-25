package sqlb

import (
	"context"
	"errors"
	"strings"
)

// Execer is an interface that wraps Exec method.
type Execer interface {
	Exec(ctx context.Context, sql string, arg []interface{}, affectedRows *int64) error
}

// InsertBuilder is a builder to build an INSERT query.
type InsertBuilder struct {
	cols         []string
	table        Table
	values       [][]interface{}
	db           Execer
	affectedRows *int64
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

// AffectedRows sets the variable to store the number of affected rows when executing the query.
func (b *InsertBuilder) AffectedRows(affectedRows *int64) *InsertBuilder {
	b.affectedRows = affectedRows
	return b
}

// SQL compiles all provided data to return an INSERT query and arguments.
func (b InsertBuilder) SQL() (string, []interface{}, error) {
	if len(b.values) == 0 {
		return "", nil, errors.New("sqlb: there must be at least one row")
	}

	sb := &strings.Builder{}
	args := argumentList{}

	_, _ = sb.WriteString("INSERT INTO ")

	b.table.Build(sb, &args)

	if len(b.cols) > 0 {
		_, _ = sb.WriteString(" (")
		for i, col := range b.cols {
			if i > 0 {
				_, _ = sb.WriteString(",")
			}
			_, _ = sb.WriteString(col)
		}
		_, _ = sb.WriteString(")")
	}

	_, _ = sb.WriteString(" VALUES ")

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

// Exec executes the INSERT query.
func (b InsertBuilder) Exec(ctx context.Context) error {
	sql, args, err := b.SQL()
	if err != nil {
		return err
	}

	return b.db.Exec(ctx, sql, args, b.affectedRows)
}
