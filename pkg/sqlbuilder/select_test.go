package sqlbuilder_test

import (
	"context"
	"encoding/json"
	"errors"
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
		require.Equal(t, "SELECT id FROM person WHERE (id = $1)", sql)
		require.Len(t, args, 1)
		require.Equal(t, []interface{}{1}, args)
	})

	t.Run("select with where and multiples conds", func(t *testing.T) {
		sql, args, err := sqlbuilder.Select("id").FromTable("person").Where(sqlbuilder.Equal("id", 1), sqlbuilder.Equal("name", "Foo")).SQL()
		require.NoError(t, err)
		require.Equal(t, "SELECT id FROM person WHERE ((id = $1) AND (name = $2))", sql)
		require.Len(t, args, 2)
		require.Equal(t, []interface{}{1, "Foo"}, args)
	})
}

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

type mockDB struct {
	err  error
	data string
	sql  string
}

func (m *mockDB) Query(ctx context.Context, query string, args []interface{}, records sqlbuilder.Recordables) error {
	m.sql = query
	_ = json.Unmarshal([]byte(m.data), records)
	return m.err
}

func (m *mockDB) QueryRow(ctx context.Context, query string, args []interface{}, record sqlbuilder.Recordable) error {
	m.sql = query
	_ = json.Unmarshal([]byte(m.data), record)
	return m.err
}

func Test_SelectBuilder_Query(t *testing.T) {
	ctx := context.Background()
	t.Run("no DB provided", func(t *testing.T) {
		f := sqlbuilder.Factory{}
		err := f.Select("id").FromTable("sample").Query(ctx, &mockRecords{})
		require.EqualError(t, err, "sqlbuilder: no DB was provided to execute the query")
	})

	t.Run("error is propagated properly", func(t *testing.T) {
		db := &mockDB{
			err: errors.New("db error"),
		}

		f := sqlbuilder.Factory{DB: db}
		err := f.Select("id").FromTable("sample").Query(ctx, &mockRecords{})
		require.EqualError(t, err, "db error")
	})

	t.Run("query is executed properly", func(t *testing.T) {
		db := &mockDB{
			data: `[{"id":1},{"id":2}]`,
		}
		f := sqlbuilder.Factory{DB: db}
		records := mockRecords{}
		err := f.Select("id").FromTable("sample").Query(ctx, &records)
		require.NoError(t, err)
		require.Equal(t, "SELECT id FROM sample", db.sql)
		require.Len(t, records, 2)
		require.Equal(t, 1, records[0].ID)
	})
}

func Test_SelectBuilder_QueryRow(t *testing.T) {
	ctx := context.Background()
	t.Run("no DB provided", func(t *testing.T) {
		f := sqlbuilder.Factory{}
		err := f.Select("id").FromTable("sample").QueryRow(ctx, &mockRecord{})
		require.EqualError(t, err, "sqlbuilder: no DB was provided to execute the query")
	})

	t.Run("error is propagated properly", func(t *testing.T) {
		db := &mockDB{
			err: errors.New("db error"),
		}

		f := sqlbuilder.Factory{DB: db}
		err := f.Select("id").FromTable("sample").QueryRow(ctx, &mockRecord{})
		require.EqualError(t, err, "db error")
	})

	t.Run("query is executed properly", func(t *testing.T) {
		db := &mockDB{
			data: `{"id":1}`,
		}
		f := sqlbuilder.Factory{DB: db}
		record := mockRecord{}
		err := f.Select("id").FromTable("sample").QueryRow(ctx, &record)
		require.NoError(t, err)
		require.Equal(t, "SELECT id FROM sample", db.sql)
		require.Equal(t, 1, record.ID)
	})
}
