package sqlbuilder

import "strings"

// Table is an interface to represent a table in SQL DB.
type Table interface {
	Name() string
}

// Column is an interface to represent a column in SQL DB.
type Column interface {
	TableName() string
	Name() string
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
		_, _ = sb.WriteString(col.TableName())
		_, _ = sb.WriteString(".")
		_, _ = sb.WriteString(col.Name())
	}

	_, _ = sb.WriteString(" FROM ")
	_, _ = sb.WriteString(b.table.Name())

	return sb.String(), nil
}
