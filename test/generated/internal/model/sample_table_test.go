package model_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bongnv/pggo/pkg/sqlbuilder"
	"github.com/bongnv/pggo/test/generated/internal/model/schema"
)

func Test_SampleTable(t *testing.T) {
	f := sqlbuilder.Factory{}
	sql, args := f.Select(schema.SampleTable.ID, schema.SampleTable.Name).From(schema.SampleTable).Build()
	require.Nil(t, args)
	require.Equal(t, "SELECT sample_table.id, sample_table.name FROM sample_table", sql)
}
