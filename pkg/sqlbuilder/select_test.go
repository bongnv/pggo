package sqlbuilder_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bongnv/pggo/pkg/sqlbuilder"
)

func Test_SelectBuilder_Build(t *testing.T) {
	t.Run("table withstring", func(t *testing.T) {
		sql, args, err := sqlbuilder.Select("id", "name").FromTable("person").SQL()
		require.NoError(t, err)
		require.Equal(t, "SELECT id, name FROM person", sql)
		require.Len(t, args, 0)
	})

	t.Run("table with type", func(t *testing.T) {
		table := sqlbuilder.BaseTable("person")
		sql, args, err := sqlbuilder.Select("id", "name").From(table).SQL()
		require.NoError(t, err)
		require.Equal(t, "SELECT id, name FROM person", sql)
		require.Len(t, args, 0)
	})

	t.Run("select with where", func(t *testing.T) {
		sql, args, err := sqlbuilder.Select("id").FromTable("person").Where(sqlbuilder.Equal("id", 1)).SQL()
		require.NoError(t, err)
		require.Equal(t, "SELECT id FROM person WHERE (id = ?)", sql)
		require.Len(t, args, 1)
		require.Equal(t, []interface{}{1}, args)
	})

	t.Run("select with where and multiples conds", func(t *testing.T) {
		sql, args, err := sqlbuilder.Select("id").FromTable("person").Where(sqlbuilder.Equal("id", 1), sqlbuilder.Equal("name", "Foo")).SQL()
		require.NoError(t, err)
		require.Equal(t, "SELECT id FROM person WHERE ((id = ?) AND (name = ?))", sql)
		require.Len(t, args, 2)
		require.Equal(t, []interface{}{1, "Foo"}, args)
	})
}
