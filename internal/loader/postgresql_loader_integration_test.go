//go:build integration
// +build integration

package loader_test

import (
	"testing"

	"github.com/bongnv/pggo/internal/loader"
	"github.com/stretchr/testify/require"
)

func Test_PostgreSQLLoader(t *testing.T) {
	l := loader.PostgreSQLLoader{
		URL: "postgres://postgres:postgres@localhost:5432/postgres",
	}
	schema, err := l.Load()
	require.NoError(t, err)
	require.Len(t, schema.Tables, 2)
	require.Contains(t, schema.Tables, "flyway_schema_history")
	require.Equal(t, "flyway_schema_history", schema.Tables["flyway_schema_history"].Name)

	sampleTable := schema.Tables["sample_table"]
	require.NotNil(t, sampleTable)
	require.Equal(t, "sample_table", sampleTable.Name)
	require.Len(t, sampleTable.Columns, 2)

	idCol := sampleTable.Columns[0]
	nameCol := sampleTable.Columns[1]

	require.Equal(t, "id", idCol.Name)
	require.False(t, idCol.Nullable)
	require.Equal(t, "uuid", idCol.DataType)

	require.Equal(t, "name", nameCol.Name)
	require.False(t, nameCol.Nullable)
	require.Equal(t, "text", nameCol.DataType)
}
