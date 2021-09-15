package sqlbuilder

import (
	"strings"
)

// Table represents a table in a database.
type Table struct {
	Name string
}

// SQL returns SQL expression of a table.
func (t Table) SQL() string {
	return t.Name
}

// ITable is an interface of a table.
type ITable interface {
	SQL() string
}

// IColumn is an interface of a column.
type IColumn interface {
	GetTable() ITable
	SQL() string
}

// Column represents a column in a table.
type Column struct {
	Table ITable
	Name  string
}

// SQL returns SQL expression of a column.
func (c Column) SQL() string {
	return c.Table.SQL() + "." + c.Name
}

func (c Column) GetTable() ITable {
	return c.Table
}

// SelectBuilder is a builder implementation of a select query.
type SelectBuilder struct {
	cols  []IColumn
	table ITable
}

func (b *SelectBuilder) From(table ITable) *SelectBuilder {
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
