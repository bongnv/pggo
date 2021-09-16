package sqlbuilder

import (
	"io"
	"strings"
)

// BaseTable represents a table in a database.
type BaseTable string

// Build adds table name to the SQL query.
func (t BaseTable) Build(sw io.StringWriter, _ ArgAppender) {
	_, _ = sw.WriteString(string(t))
}

func (t BaseTable) tableOnly() {}

// Builder is the interface that wraps the Build method.
type Builder interface {
	// Build contructs SQL and update arguments if necessary.
	Build(sw io.StringWriter, args ArgAppender)
}

// Table is an interface of a table.
type Table interface {
	Builder
	tableOnly()
}

// SelectBuilder is a builder implementation of a select query.
type SelectBuilder struct {
	cols      []string
	table     Table
	conds     []Condition
	arguments *argumentList
}

// FromSQL sets the FROM clause for the query.
func (b *SelectBuilder) FromSQL(table string) *SelectBuilder {
	b.table = BaseTable(table)
	return b
}

// From sets the FROM clause from the given table for the query.
func (b *SelectBuilder) From(table Table) *SelectBuilder {
	b.table = table
	return b
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

	if b.table != nil {
		_, _ = sb.WriteString(" FROM ")
		b.table.Build(sb, b.arguments)
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

func (l *argumentList) Append(values ...interface{}) {
	l.args = append(l.args, values...)
}
