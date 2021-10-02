package sqlb

import (
	"context"
	"io"
	"strconv"
	"strings"
)

// BaseTable represents a table in a database.
type BaseTable string

// Build adds table name to the SQL query.
func (t BaseTable) Build(sw io.StringWriter, _ Placeholders) error {
	_, _ = sw.WriteString(string(t))
	return nil
}

func (t BaseTable) tableOnly() {}

// Builder is the interface that wraps the Build method.
type Builder interface {
	// Build contructs SQL and update arguments if necessary.
	Build(sw io.StringWriter, args Placeholders) error
}

// Table is an interface of a table.
type Table interface {
	Builder
	tableOnly()
}

// Queryer is an interface that wraps Query functinos.
type Queryer interface {
	Query(ctx context.Context, query string, args []interface{}, records EntityList) error
	QueryRow(ctx context.Context, query string, args []interface{}, record Entity) error
}

// SelectBuilder is a builder implementation of a select query.
type SelectBuilder struct {
	cols  []string
	db    Queryer
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
func (b SelectBuilder) Build(sb io.StringWriter, aa Placeholders) error {
	_, _ = sb.WriteString("SELECT ")

	for i, col := range b.cols {
		if i > 0 {
			_, _ = sb.WriteString(", ")
		}
		_, _ = sb.WriteString(col)
	}

	if b.from != nil {
		if err := b.from.Build(sb, aa); err != nil {
			return err
		}
	}

	if b.where != nil {
		if err := b.where.Build(sb, aa); err != nil {
			return err
		}
	}

	return nil
}

// SQL compiles all provided data to return a SELECT query and arguments.
func (b SelectBuilder) SQL() (string, []interface{}, error) {
	sb := &strings.Builder{}
	args := argumentList{}
	if err := b.Build(sb, &args); err != nil {
		return "", nil, err
	}

	return sb.String(), args, nil
}

// Query sends the query to DB and parses results to the given records.
func (b SelectBuilder) Query(ctx context.Context, records EntityList) error {
	sql, args, err := b.SQL()
	if err != nil {
		return err
	}

	return b.db.Query(ctx, sql, args, records)
}

// QueryRow sends the query to DB and parses results to the given record.
// If no rows were found it returns ErrNoRows. If multiple rows are returned it
// ignores all but the first.
func (b SelectBuilder) QueryRow(ctx context.Context, record Entity) error {
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

func (c fromClause) Build(sb io.StringWriter, aa Placeholders) error {
	_, _ = sb.WriteString(" FROM ")
	return c.table.Build(sb, aa)
}
