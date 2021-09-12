package generator

import (
	"bytes"
	"fmt"

	"github.com/bongnv/pggo/internal/template"
)

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
	Tables map[string]*Table
}

// SchemaLoader is an interface that wraps Load method.
type SchemaLoader interface {
	// Load loads DB schema. It returns error if there is any error while loading the schema.
	Load() (*Schema, error)
}

// Writer is an interface that wraps Write method.
type Writer interface {
	// Write writes content into a file.
	Write(fileName string, content []byte) error
}

type templateData struct {
	PackageName string
	Table       *Table
}

// Generator is an implementation to generate Go code from schema.
type Generator struct {
	SchemaLoader SchemaLoader
	Table        string
	Writer       Writer
}

// Generate generates Go code from DB schema.
func (g Generator) Generate() error {
	schema, err := g.SchemaLoader.Load()
	if err != nil {
		return err
	}

	table := schema.Tables[g.Table]
	if table == nil {
		return fmt.Errorf("generator: %s couldn't be found in the schema", g.Table)
	}

	data := &templateData{
		PackageName: "model",
		Table:       schema.Tables[g.Table],
	}
	buf := &bytes.Buffer{}
	if err := template.Execute(buf, "table_model.tmpl", data); err != nil {
		return err
	}

	return g.Writer.Write(g.Table+".pggo.go", buf.Bytes())
}
