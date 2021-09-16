package sqlbuilder

import (
	"strings"
)

// BaseTable represents a table in a database.
type BaseTable string

// SQL returns SQL expression of a table.
func (t BaseTable) SQL() string {
	return string(t)
}

// Table is an interface of a table.
type Table interface {
	SQL() string
}

// SelectBuilder is a builder implementation of a select query.
type SelectBuilder struct {
	cols      []string
	table     string
	conds     []Condition
	arguments *argumentList
}

// FromSQL sets the FROM clause for the query.
func (b *SelectBuilder) FromSQL(table string) *SelectBuilder {
	b.table = table
	return b
}

// From sets the FROM clause from the given table for the query.
func (b *SelectBuilder) From(table Table) *SelectBuilder {
	return b.FromSQL(table.SQL())
}

// Where sets the WHERE clause for the query.
func (b *SelectBuilder) Where(conds ...Condition) *SelectBuilder {
	b.conds = conds
	return b
}

// Build compiles all provided data to return a SELECT query and arguments.
func (b SelectBuilder) Build() (string, []interface{}, error) {
	sb := &strings.Builder{}
	_, _ = sb.WriteString("SELECT ")

	for i, col := range b.cols {
		if i > 0 {
			_, _ = sb.WriteString(", ")
		}
		_, _ = sb.WriteString(col)
	}

	if b.table != "" {
		_, _ = sb.WriteString(" FROM ")
		_, _ = sb.WriteString(b.table)
	}

	if len(b.conds) > 0 {
		_, _ = sb.WriteString(" WHERE ")
		for i, cond := range b.conds {
			if i > 0 {
				_, _ = sb.WriteString(" AND ")
			}
			cond(sb, b.arguments)
		}
	}

	return sb.String(), b.arguments.args, nil
}

type argumentList struct {
	args []interface{}
}

func (l *argumentList) Append(v interface{}) {
	l.args = append(l.args, v)
}
