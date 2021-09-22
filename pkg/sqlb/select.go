package sqlb

import (
	"context"
	"io"
	"strconv"
	"strings"
)

// Select starts a new SELECT query.
func Select(cols ...string) *SelectBuilder {
	return Factory{}.Select(cols...)
}

// BaseTable represents a table in a database.
type BaseTable string

// Build adds table name to the SQL query.
func (t BaseTable) Build(sw io.StringWriter, _ Placeholders) {
	_, _ = sw.WriteString(string(t))
}

func (t BaseTable) tableOnly() {}

// Builder is the interface that wraps the Build method.
type Builder interface {
	// Build contructs SQL and update arguments if necessary.
	Build(sw io.StringWriter, args Placeholders)
}

// Table is an interface of a table.
type Table interface {
	Builder
	tableOnly()
}

// SelectBuilder is a builder implementation of a select query.
type SelectBuilder struct {
	cols  []string
	db    DB
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
func (b SelectBuilder) Build(sb io.StringWriter, aa Placeholders) {
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
	args := argumentList{}
	b.Build(sb, &args)
	return sb.String(), args, nil
}

// Query sends the query to DB and parses results to the given records.
func (b SelectBuilder) Query(ctx context.Context, records Recordables) error {
	sql, args, err := b.SQL()
	if err != nil {
		return err
	}

	return b.db.Query(ctx, sql, args, records)
}

// QueryRow sends the query to DB and parses results to the given record.
// If no rows were found it returns ErrNoRows. If multiple rows are returned it
// ignores all but the first.
func (b SelectBuilder) QueryRow(ctx context.Context, record Recordable) error {
	sql, args, err := b.SQL()
	if err != nil {
		return err
	}

	return b.db.QueryRow(ctx, sql, args, record)
}

type argumentList []interface{}

func (l *argumentList) Append(values ...interface{}) string {
	*l = append(*l, values...)
	return "$" + strconv.Itoa(len(*l))
}

type fromClause struct {
	table Table
}

func (c fromClause) Build(sb io.StringWriter, aa Placeholders) {
	_, _ = sb.WriteString(" FROM ")
	c.table.Build(sb, aa)
}
