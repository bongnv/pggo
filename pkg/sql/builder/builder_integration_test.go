//go:build integration
// +build integration

package builder_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/bongnv/pggo/pkg/sql/builder"
	sqlb "github.com/bongnv/pggo/pkg/sqlb"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/require"
)

func Test_QueryRow(t *testing.T) {
	ctx := context.Background()
	conn, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:5432/postgres")
	require.NoError(t, err)
	defer conn.Close()

	t.Run("happy-path", func(t *testing.T) {
		r := &mockRecord{}
		err = builder.With(conn).
			Select("id", "name").
			FromTable("sample_table").
			Where(sqlb.Equal("id", 1)).
			QueryRow(ctx, r)
		require.NoError(t, err)
		require.Equal(t, 1, r.ID)
		require.Equal(t, "One", r.Name)
	})

	t.Run("missing-name", func(t *testing.T) {
		r := &mockRecord{}
		err := builder.With(conn).
			Select("id", "name as notfound").
			FromTable("sample_table").
			Where(sqlb.Equal("id", 1)).
			QueryRow(ctx, r)
		require.EqualError(t, err, "notfound doesn't exist")
	})

	t.Run("no-record", func(t *testing.T) {
		r := &mockRecord{}
		err = builder.With(conn).
			Select("id", "name").
			FromTable("sample_table").
			Where(sqlb.Equal("id", 2)).
			QueryRow(ctx, r)
		require.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("no-record", func(t *testing.T) {
		r := &mockRecord{}
		err = builder.With(conn).
			Select("id", "name").
			FromTable("nonexist_table").
			QueryRow(ctx, r)
		require.EqualError(t, err, "ERROR: relation \"nonexist_table\" does not exist (SQLSTATE 42P01)")
	})

	t.Run("wrong column type", func(t *testing.T) {
		r := &mockRecord{}
		err = builder.With(conn).
			Select("name as number").
			FromTable("sample_table").
			Where(sqlb.Equal("id", 1)).
			QueryRow(ctx, r)
		require.EqualError(t, err, "sql: Scan error on column index 0, name \"number\": converting driver.Value type string (\"One\") to a int: invalid syntax")
	})
}

func Test_Query(t *testing.T) {
	ctx := context.Background()
	conn, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:5432/postgres")
	require.NoError(t, err)
	defer conn.Close()

	t.Run("happy-path", func(t *testing.T) {
		records := mockRecords{}
		err = builder.With(conn).
			Select("id", "name").
			FromTable("sample_table").
			Where(sqlb.Equal("id", 1)).
			Query(ctx, &records)
		require.NoError(t, err)
		require.Len(t, records, 1)
		require.Equal(t, 1, records[0].ID)
		require.Equal(t, "One", records[0].Name)
	})

	t.Run("missing-name", func(t *testing.T) {
		records := mockRecords{}
		err := builder.With(conn).
			Select("id", "name as notfound").
			FromTable("sample_table").
			Where(sqlb.Equal("id", 1)).
			Query(ctx, &records)
		require.EqualError(t, err, "notfound doesn't exist")
	})

	t.Run("no-record", func(t *testing.T) {
		records := mockRecords{}
		err = builder.With(conn).
			Select("id", "name").
			FromTable("nonexist_table").
			Query(ctx, &records)
		require.EqualError(t, err, "ERROR: relation \"nonexist_table\" does not exist (SQLSTATE 42P01)")
	})

	t.Run("wrong column type", func(t *testing.T) {
		records := mockRecords{}
		err = builder.With(conn).
			Select("name as number").
			FromTable("sample_table").
			Where(sqlb.Equal("id", 1)).
			Query(ctx, &records)
		require.EqualError(t, err, "sql: Scan error on column index 0, name \"number\": converting driver.Value type string (\"One\") to a int: invalid syntax")
	})
}

func Test_Exec(t *testing.T) {
	ctx := context.Background()
	conn, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:5432/postgres")
	require.NoError(t, err)
	defer conn.Close()

	teardown := func() {
		_, err := conn.Exec("delete from sample_table where id IN (2)")
		require.NoError(t, err)
	}

	t.Run("table not found", func(t *testing.T) {
		err := builder.With(conn).InsertTable("not_found").Values(1, "Joe").Exec(ctx)
		require.EqualError(t, err, "ERROR: relation \"not_found\" does not exist (SQLSTATE 42P01)")
	})

	t.Run("insert successfully", func(t *testing.T) {
		defer teardown()
		var affectedRows int64
		err := builder.With(conn).
			InsertTable("sample_table").
			Columns("id", "name").
			Values(2, "Joe Two").
			AffectedRows(&affectedRows).
			Exec(ctx)
		require.NoError(t, err)
		require.EqualValues(t, 1, affectedRows)
	})
}
