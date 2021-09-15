package sqlbuilder

// Factory is a factory implementation for SQL builders.
type Factory struct{}

// Select creates a new builder for a SELECT query.
func (f Factory) Select(cols ...Column) *SelectBuilder {
	return &SelectBuilder{
		cols: cols,
	}
}
