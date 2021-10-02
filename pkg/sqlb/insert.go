package sqlb

import (
	"context"
	"errors"
	"io"
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
	values       []Builder
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
	b.values = append(b.values, groupPlaceholder{
		values: values,
	})
	return b
}

// AffectedRows sets the variable to store the number of affected rows when executing the query.
func (b *InsertBuilder) AffectedRows(affectedRows *int64) *InsertBuilder {
	b.affectedRows = affectedRows
	return b
}

// Entities adds a single row or multiple rows to the query via Entity objects.
func (b *InsertBuilder) Entities(entities ...Entity) *InsertBuilder {
	for _, e := range entities {
		var valBuilder builderFn = func(sw io.StringWriter, args Placeholders) error {
			values, err := e.GetValues(b.cols)
			if err != nil {
				return err
			}

			return groupPlaceholder{
				values: values,
			}.Build(sw, args)
		}
		b.values = append(b.values, valBuilder)
	}
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

	if err := b.table.Build(sb, &args); err != nil {
		return "", nil, err
	}

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

	for i, valBuilder := range b.values {
		if i > 0 {
			_, _ = sb.WriteString(",")
		}
		if err := valBuilder.Build(sb, &args); err != nil {
			return "", nil, err
		}
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

type builderFn func(sw io.StringWriter, args Placeholders) error

func (f builderFn) Build(sw io.StringWriter, args Placeholders) error {
	return f(sw, args)
}
