//go:build integration
// +build integration

package sqlb_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/bongnv/pggo/pkg/sql/sqlb"
	"github.com/bongnv/pggo/pkg/sqlbuilder"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/require"
)

type mockRecord struct {
	ID   int
	Name string
}

func (m *mockRecord) GetPointer(col string) interface{} {
	switch col {
	case "id":
		return &m.ID
	case "name":
		return &m.Name
	default:
		return nil
	}
}

type mockRecords []*mockRecord

func (m mockRecords) New() sqlbuilder.Recordable {
	return &mockRecord{}
}

func (m *mockRecords) Append(r sqlbuilder.Recordable) {
	*m = append(*m, r.(*mockRecord))
}

func Test_QueryRow(t *testing.T) {
	ctx := context.Background()
	conn, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:5432/postgres")
	require.NoError(t, err)
	defer conn.Close()

	t.Run("happy-path", func(t *testing.T) {
		r := &mockRecord{}
		err = sqlb.With(conn).
			Select("id", "name").
			FromTable("sample_table").
			Where(sqlbuilder.Equal("id", 1)).
			QueryRow(ctx, r)
		require.NoError(t, err)
		require.Equal(t, 1, r.ID)
		require.Equal(t, "One", r.Name)
	})

	t.Run("missing-name", func(t *testing.T) {
		r := &mockRecord{}
		err := sqlb.With(conn).
			Select("id", "name as notfound").
			FromTable("sample_table").
			Where(sqlbuilder.Equal("id", 1)).
			QueryRow(ctx, r)
		require.EqualError(t, err, "sqlb: notfound is not found")
	})

	t.Run("no-record", func(t *testing.T) {
		r := &mockRecord{}
		err = sqlb.With(conn).
			Select("id", "name").
			FromTable("sample_table").
			Where(sqlbuilder.Equal("id", 2)).
			QueryRow(ctx, r)
		require.Equal(t, sql.ErrNoRows, err)
	})
}

func Test_Query(t *testing.T) {
	ctx := context.Background()
	conn, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:5432/postgres")
	require.NoError(t, err)
	defer conn.Close()

	t.Run("happy-path", func(t *testing.T) {
		records := mockRecords{}
		err = sqlb.With(conn).
			Select("id", "name").
			FromTable("sample_table").
			Where(sqlbuilder.Equal("id", 1)).
			Query(ctx, &records)
		require.NoError(t, err)
		require.Len(t, records, 1)
		require.Equal(t, 1, records[0].ID)
		require.Equal(t, "One", records[0].Name)
	})

	t.Run("missing-name", func(t *testing.T) {
		records := mockRecords{}
		err := sqlb.With(conn).
			Select("id", "name as notfound").
			FromTable("sample_table").
			Where(sqlbuilder.Equal("id", 1)).
			Query(ctx, &records)
		require.EqualError(t, err, "sqlb: notfound is not found")
	})
}
