package sqlb_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bongnv/pggo/pkg/sqlb"
)

func Test_Insert_SQL(t *testing.T) {
	t.Run("via table name", func(t *testing.T) {
		sql, args, err := sqlb.InsertTable("person").Columns("id", "name").Values(1, "Joe").SQL()
		require.NoError(t, err)
		require.Equal(t, "INSERT INTO person (id,name) VALUES ($1,$2)", sql)
		require.Equal(t, []interface{}{1, "Joe"}, args)
	})

	t.Run("via table", func(t *testing.T) {
		sql, args, err := sqlb.Insert(sqlb.BaseTable("person")).Columns("id", "name").Values(1, "Joe").SQL()
		require.NoError(t, err)
		require.Equal(t, "INSERT INTO person (id,name) VALUES ($1,$2)", sql)
		require.Equal(t, []interface{}{1, "Joe"}, args)
	})

	t.Run("insert with multiple rows", func(t *testing.T) {
		sql, args, err := sqlb.InsertTable("person").Columns("id", "name").Values(1, "Joe").Values(2, "Two").SQL()
		require.NoError(t, err)
		require.Equal(t, "INSERT INTO person (id,name) VALUES ($1,$2),($3,$4)", sql)
		require.Equal(t, []interface{}{1, "Joe", 2, "Two"}, args)
	})

	t.Run("zero row of values", func(t *testing.T) {
		sql, args, err := sqlb.InsertTable("person").Columns("id", "name").SQL()
		require.EqualError(t, err, "sqlb: there must be at least one row")
		require.Empty(t, sql)
		require.Empty(t, args)
	})
}

type mockExecer struct {
	called       bool
	err          error
	sql          string
	args         []interface{}
	affectedRows int64
}

func (m *mockExecer) Exec(_ context.Context, sql string, args []interface{}, affectedRows *int64) error {
	m.called = true
	m.sql = sql
	m.args = args
	if affectedRows != nil {
		*affectedRows = m.affectedRows
	}
	return m.err
}

func Test_Insert_Exec(t *testing.T) {
	ctx := context.Background()

	t.Run("without db", func(t *testing.T) {
		err := sqlb.InsertTable("person").Values(t, "Joe").Exec(ctx)
		require.EqualError(t, err, "sqlb: no DB was provided to execute the query")
	})

	t.Run("error when creating SQL", func(t *testing.T) {
		db := &mockExecer{}
		err := sqlb.MakeInsertBuilder(db, "person").Columns("id", "name").Exec(ctx)
		require.EqualError(t, err, "sqlb: there must be at least one row")
		require.False(t, db.called)
	})

	t.Run("run with error from DB", func(t *testing.T) {
		db := &mockExecer{
			err: errors.New("remote error"),
		}
		err := sqlb.MakeInsertBuilder(db, "person").Columns("id", "name").Values(1, "Joe").Exec(ctx)
		require.EqualError(t, err, "remote error")
		require.True(t, db.called)
		require.Equal(t, "INSERT INTO person (id,name) VALUES ($1,$2)", db.sql)
		require.Equal(t, []interface{}{1, "Joe"}, db.args)
	})

	t.Run("run with no error", func(t *testing.T) {
		db := &mockExecer{
			affectedRows: 1,
		}
		var affectedRows int64
		err := sqlb.MakeInsertBuilder(db, "person").Values(1, "Joe").AffectedRows(&affectedRows).Exec(ctx)
		require.NoError(t, err)
		require.True(t, db.called)
		require.Equal(t, "INSERT INTO person VALUES ($1,$2)", db.sql)
		require.Equal(t, []interface{}{1, "Joe"}, db.args)
		require.EqualValues(t, 1, affectedRows)
	})
}
