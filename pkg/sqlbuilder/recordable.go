package sqlbuilder

// Recordable stores value from a DB row.
type Recordable interface {
	// GetPointer returns the pointer for storing data from a given column. It returns nil if the column doesn't match.
	GetPointer(col string) interface{}
}

// Recordables represents a list of Recordable.
type Recordables interface {
	New() Recordable
	Append(r Recordable)
}
