package sqlbuilder_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bongnv/pggo/pkg/sqlbuilder"
)

func Test_Insert_SQL(t *testing.T) {
	t.Run("via table name", func(t *testing.T) {
		sql, args, err := sqlbuilder.InsertTable("person").Columns("id", "name").Values(1, "Joe").SQL()
		require.NoError(t, err)
		require.Equal(t, "INSERT INTO person (id,name) VALUES ($1,$2)", sql)
		require.Equal(t, []interface{}{1, "Joe"}, args)
	})

	t.Run("via table", func(t *testing.T) {
		sql, args, err := sqlbuilder.Insert(sqlbuilder.BaseTable("person")).Columns("id", "name").Values(1, "Joe").SQL()
		require.NoError(t, err)
		require.Equal(t, "INSERT INTO person (id,name) VALUES ($1,$2)", sql)
		require.Equal(t, []interface{}{1, "Joe"}, args)
	})

	t.Run("insert with multiple rows", func(t *testing.T) {
		sql, args, err := sqlbuilder.InsertTable("person").Columns("id", "name").Values(1, "Joe").Values(2, "Two").SQL()
		require.NoError(t, err)
		require.Equal(t, "INSERT INTO person (id,name) VALUES ($1,$2),($3,$4)", sql)
		require.Equal(t, []interface{}{1, "Joe", 2, "Two"}, args)
	})

	t.Run("zero row of values", func(t *testing.T) {
		sql, args, err := sqlbuilder.InsertTable("person").Columns("id", "name").SQL()
		require.EqualError(t, err, "sqlbuilder: there must be at least one row")
		require.Empty(t, sql)
		require.Empty(t, args)
	})
}
