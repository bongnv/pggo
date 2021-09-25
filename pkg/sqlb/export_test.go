package sqlb

// Bridge package to expose internals for testing.

// ArgumentList is exported for testing.
type ArgumentList = argumentList

// MakeInsertBuilder is exported for testing.
func MakeInsertBuilder(db Execer, table string) *InsertBuilder {
	return &InsertBuilder{
		db:    db,
		table: BaseTable(table),
	}
}

// MakeSelectBuilder is exported for testing.
func MakeSelectBuilder(db Queryer, cols ...string) *SelectBuilder {
	return &SelectBuilder{
		db:   db,
		cols: cols,
	}
}
