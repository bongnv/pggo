package sqlschema

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

// Column represents a column in a table.
type Column struct {
	Table ITable
	Name  string
}

// SQL returns SQL expression of a column.
func (c Column) SQL() string {
	return c.Table.SQL() + "." + c.Name
}
