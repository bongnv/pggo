package sqlbuilder

import (
	"strings"
)

// BaseTable represents a table in a database.
type BaseTable struct {
	Name string
}

// SQL returns SQL expression of a table.
func (t BaseTable) SQL() string {
	return t.Name
}

// Table is an interface of a table.
type Table interface {
	SQL() string
}

// Column is an interface of a column.
type Column interface {
	SQL() string
}

// Column represents a column in a table.
type BaseColumn struct {
	Table Table
	Name  string
}

// SQL returns SQL expression of a column.
func (c BaseColumn) SQL() string {
	return c.Table.SQL() + "." + c.Name
}

// SelectBuilder is a builder implementation of a select query.
type SelectBuilder struct {
	cols  []Column
	table Table
}

func (b *SelectBuilder) From(table Table) *SelectBuilder {
	b.table = table
	return b
}

// Build compiles all provided data to return a SELECT query and arguments.
func (b SelectBuilder) Build() (string, []interface{}) {
	sb := &strings.Builder{}
	_, _ = sb.WriteString("SELECT ")

	for i, col := range b.cols {
		if i > 0 {
			_, _ = sb.WriteString(", ")
		}
		_, _ = sb.WriteString(col.SQL())
	}

	_, _ = sb.WriteString(" FROM ")
	_, _ = sb.WriteString(b.table.SQL())

	return sb.String(), nil
}
