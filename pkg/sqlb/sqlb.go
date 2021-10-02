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

// Entity represents an entity usually a row in a table.
type Entity interface {
	// GetPointers returns the pointer for storing data from a given column. It returns nil if the column doesn't match.
	GetPointers(cols []string) ([]interface{}, error)
	// GetValues returns the value given a column. It returns nil if the column doesn't match.
	GetValues(cols []string) ([]interface{}, error)
}

// EntityList represents a list of entities.
type EntityList interface {
	// New is a factory function to create a new entity.
	New() Entity
	// Append adds a new entity into the list.
	Append(e Entity)
}
