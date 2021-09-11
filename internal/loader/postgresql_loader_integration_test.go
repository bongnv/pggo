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
	require.Equal(t, "flyway_schema_history", schema.Tables[0].Name)
	require.Equal(t, "sample_table", schema.Tables[1].Name)
}
