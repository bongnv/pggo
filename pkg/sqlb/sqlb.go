package sqlb

// Insert starts a new INSERT query.
func Insert(table Table) *InsertBuilder {
	return DefaultFactory.Insert(table)
}

// Insert starts a new INSERT query with a table name.
func InsertTable(tableName string) *InsertBuilder {
	return DefaultFactory.InsertTable(tableName)
}

// Select starts a new SELECT query.
func Select(cols ...string) *SelectBuilder {
	return DefaultFactory.Select(cols...)
}
