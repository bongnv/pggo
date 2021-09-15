package sqlbuilder_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bongnv/pggo/pkg/sqlbuilder"
)

func Test_SelectBuilder_Build_Happy(t *testing.T) {
	f := sqlbuilder.Factory{}
	t.Run("table withstring", func(t *testing.T) {
		sql, args, err := f.Select("id", "name").FromSQL("person").Build()
		require.NoError(t, err)
		require.Equal(t, "SELECT id, name FROM person", sql)
		require.Len(t, args, 0)
	})

	t.Run("table with type", func(t *testing.T) {
		table := sqlbuilder.BaseTable("person")
		sql, args, err := f.Select("id", "name").From(table).Build()
		require.NoError(t, err)
		require.Equal(t, "SELECT id, name FROM person", sql)
		require.Len(t, args, 0)
	})
}
