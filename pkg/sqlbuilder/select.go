package sqlbuilder

import (
	"errors"
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
	cols  []string
	table Table
	err   error
}

func (b *SelectBuilder) From(table interface{}) *SelectBuilder {
	switch v := table.(type) {
	case string:
		b.table = BaseTable(v)
	case Table:
		b.table = v
	default:
		b.err = errors.New("sqlbuilder: the given table has invalid type")
	}
	return b
}

// Build compiles all provided data to return a SELECT query and arguments.
func (b SelectBuilder) Build() (string, []interface{}, error) {
	if b.err != nil {
		return "", nil, b.err
	}

	sb := &strings.Builder{}
	_, _ = sb.WriteString("SELECT ")

	for i, col := range b.cols {
		if i > 0 {
			_, _ = sb.WriteString(", ")
		}
		_, _ = sb.WriteString(col)
	}

	_, _ = sb.WriteString(" FROM ")
	_, _ = sb.WriteString(b.table.SQL())

	return sb.String(), nil, nil
}
