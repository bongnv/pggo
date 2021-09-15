package generator_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bongnv/pggo/internal/generator"
)

type mockSchemaLoader struct {
	Schema *generator.Schema
	Err    error
}

func (l mockSchemaLoader) Load() (*generator.Schema, error) {
	return l.Schema, l.Err
}

type mockWriter struct {
	bytes.Buffer
}

func (w *mockWriter) Write(fileName string, content []byte) error {
	if _, err := fmt.Fprintln(&w.Buffer, fileName); err != nil {
		return err
	}

	_, err := w.Buffer.Write(content)
	return err
}

func Test_Generator_happy(t *testing.T) {
	loader := &mockSchemaLoader{
		Schema: &generator.Schema{
			Tables: map[string]*generator.Table{
				"mock_table": {
					Name: "mock_table",
					Columns: []*generator.Column{
						{
							Name:     "id",
							DataType: "uuid",
						},
						{
							Name:     "name",
							DataType: "text",
						},
					},
				},
			},
		},
	}
	writer := &mockWriter{}
	g := &generator.Generator{
		SchemaLoader: loader,
		Table:        "mock_table",
		Writer:       writer,
	}
	require.NoError(t, g.Generate())
	require.Equal(t, `mock_table.pggo.go
package model

import (
	"github.com/google/uuid"
)

// SampleTable represents mock_table table.
type SampleTable struct {
	ID   uuid.UUID
	Name string
}
schema/mock_table.pggo.go
package schema

import "github.com/bongnv/pggo/pkg/sqlbuilder"

// SampleTable defines the schema of mock_table.
var SampleTable = struct {
	sqlbuilder.Table
	ID   sqlbuilder.Column
	Name sqlbuilder.Column
}{
	Table: sqlbuilder.Table{
		Name: "mock_table",
	},
}

func init() {
	SampleTable.ID = sqlbuilder.Column{
		Table: SampleTable,
		Name:  "id",
	}

	SampleTable.Name = sqlbuilder.Column{
		Table: SampleTable,
		Name:  "name",
	}
}
`,
		writer.String())
}

func Test_Generator_table_not_found(t *testing.T) {
	loader := &mockSchemaLoader{
		Schema: &generator.Schema{
			Tables: map[string]*generator.Table{
				"mock_table": {
					Name: "mock_table",
				},
			},
		},
	}
	writer := &mockWriter{}
	g := &generator.Generator{
		SchemaLoader: loader,
		Table:        "mock_table_not_found",
		Writer:       writer,
	}
	require.EqualError(t, g.Generate(), "generator: mock_table_not_found couldn't be found in the schema")
}

func Test_Generator_schema_error(t *testing.T) {
	loader := &mockSchemaLoader{
		Err: errors.New("random error"),
	}
	writer := &mockWriter{}
	g := &generator.Generator{
		SchemaLoader: loader,
		Table:        "mock_table",
		Writer:       writer,
	}
	require.EqualError(t, g.Generate(), "random error")
}
