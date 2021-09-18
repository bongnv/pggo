package model_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bongnv/pggo/pkg/sqlbuilder"
	"github.com/bongnv/pggo/test/generated/internal/model/schema"
)

func Test_SampleTable(t *testing.T) {
	sql, args, err := sqlbuilder.Select(schema.SampleTable.ID, schema.SampleTable.Name).From(schema.SampleTable).SQL()

	require.NoError(t, err)
	require.Empty(t, args)
	require.Equal(t, "SELECT id, name FROM sample_table", sql)
}
