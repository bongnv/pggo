package generator

import "log"

// Column represents a column in a table.
type Column struct {
	Name     string
	Nullable bool
	DataType string
}

// Table represents a table in a schema.
type Table struct {
	Name    string
	Columns []*Column
}

// Scheme represents a DB schema.
type Schema struct {
	Tables []*Table
}

// SchemaLoader is an interface that wraps Load method.
type SchemaLoader interface {
	// Load loads DB schema. It returns error if there is any error while loading the schema.
	Load() (*Schema, error)
}

// Generator is an implementation to generate Go code from schema.
type Generator struct {
	SchemaLoader SchemaLoader
}

func (g Generator) Generate() error {
	schema, err := g.SchemaLoader.Load()
	if err != nil {
		return err
	}

	log.Println("Loaded schema:", schema)
	return nil
}
