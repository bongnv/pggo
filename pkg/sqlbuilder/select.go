package sqlbuilder

import (
	"io"
	"strings"
)

// Select starts a new SELECT query.
func Select(cols ...string) *SelectBuilder {
	return &SelectBuilder{
		cols: cols,
	}
}

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
	cols  []string
	from  Builder
	where Builder
}

// FromTable sets the FROM clause for the query with the table is provided with a string.
func (b *SelectBuilder) FromTable(table string) *SelectBuilder {
	return b.From(BaseTable(table))
}

// From sets the FROM clause from the given table for the query.
func (b *SelectBuilder) From(table Table) *SelectBuilder {
	b.from = fromClause{
		table: table,
	}
	return b
}

// Where sets the WHERE clause for the query.
func (b *SelectBuilder) Where(conds ...Condition) *SelectBuilder {
	b.where = whereClause{
		cond: And(conds...),
	}
	return b
}

// Build builds the SELECT query.
func (b SelectBuilder) Build(sb io.StringWriter, aa ArgAppender) {
	_, _ = sb.WriteString("SELECT ")

	for i, col := range b.cols {
		if i > 0 {
			_, _ = sb.WriteString(", ")
		}
		_, _ = sb.WriteString(col)
	}

	if b.from != nil {
		b.from.Build(sb, aa)
	}

	if b.where != nil {
		b.where.Build(sb, aa)
	}
}

// SQL compiles all provided data to return a SELECT query and arguments.
func (b SelectBuilder) SQL() (string, []interface{}, error) {
	sb := &strings.Builder{}
	args := &argumentList{}
	b.Build(sb, args)
	return sb.String(), args.Args, nil
}

type argumentList struct {
	Args []interface{}
}

func (l *argumentList) Append(values ...interface{}) {
	l.Args = append(l.Args, values...)
}

type fromClause struct {
	table Table
}

func (c fromClause) Build(sb io.StringWriter, aa ArgAppender) {
	_, _ = sb.WriteString(" FROM ")
	c.table.Build(sb, aa)
}
