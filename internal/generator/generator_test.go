package generator_test

import (
	"bytes"
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

func Test_Generator(t *testing.T) {
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
		Table:        "mock_table",
		Writer:       writer,
	}
	require.NoError(t, g.Generate())
	require.Equal(t, "mock_table.pggo.go\npackage model\n\nimport (\n\t\"github.com/google/uuid\"\n)\n\n// SampleTable represents mock_table table.\ntype SampleTable struct {\n\tID   uuid.UUID\n\tName string\n}\n", writer.String())
}