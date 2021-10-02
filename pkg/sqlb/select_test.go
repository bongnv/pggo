package sqlb_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bongnv/pggo/pkg/sqlb"
)

func Test_SelectBuilder_SQL(t *testing.T) {
	t.Run("table withstring", func(t *testing.T) {
		sql, args, err := sqlb.Select("id", "name").FromTable("person").SQL()
		require.NoError(t, err)
		require.Equal(t, "SELECT id, name FROM person", sql)
		require.Len(t, args, 0)
	})

	t.Run("table with type", func(t *testing.T) {
		table := sqlb.BaseTable("person")
		sql, args, err := sqlb.Select("id", "name").From(table).SQL()
		require.NoError(t, err)
		require.Equal(t, "SELECT id, name FROM person", sql)
		require.Len(t, args, 0)
	})

	t.Run("select with where", func(t *testing.T) {
		sql, args, err := sqlb.Select("id").FromTable("person").Where(sqlb.Equal("id", 1)).SQL()
		require.NoError(t, err)
		require.Equal(t, "SELECT id FROM person WHERE (id = $1)", sql)
		require.Len(t, args, 1)
		require.Equal(t, []interface{}{1}, args)
	})

	t.Run("select with error where", func(t *testing.T) {
		_, _, err := sqlb.Select("id").Where(sqlb.In("id")).SQL()
		require.EqualError(t, err, "values list must not be empty")
	})

	t.Run("select with where and multiples conds", func(t *testing.T) {
		sql, args, err := sqlb.Select("id").FromTable("person").Where(sqlb.Equal("id", 1), sqlb.Equal("name", "Foo")).SQL()
		require.NoError(t, err)
		require.Equal(t, "SELECT id FROM person WHERE ((id = $1) AND (name = $2))", sql)
		require.Len(t, args, 2)
		require.Equal(t, []interface{}{1, "Foo"}, args)
	})

	t.Run("select from error table", func(t *testing.T) {
		_, _, err := sqlb.Select("id").From(tableWithErr{err: errors.New("random error")}).SQL()
		require.EqualError(t, err, "random error")
	})
}

type mockRecord struct {
	ID   int
	Name string
}

func (m *mockRecord) GetPointers(cols []string) ([]interface{}, error) {
	if len(cols) == 0 {
		return []interface{}{&m.ID, &m.Name}, nil
	}

	pointers := make([]interface{}, len(cols))
	for i, col := range cols {
		switch col {
		case "id":
			pointers[i] = &m.ID
		case "name":
			pointers[i] = &m.Name
		default:
			return nil, fmt.Errorf("%s couldn't be found", col)
		}
	}

	return pointers, nil
}

func (m *mockRecord) GetValues(cols []string) ([]interface{}, error) {
	if len(cols) == 0 {
		return []interface{}{m.ID, m.Name}, nil
	}

	values := make([]interface{}, len(cols))
	for i, col := range cols {
		switch col {
		case "id":
			values[i] = m.ID
		case "name":
			values[i] = m.Name
		default:
			return nil, fmt.Errorf("%s couldn't be found", col)
		}
	}

	return values, nil
}

type mockRecords []*mockRecord

func (m mockRecords) New() sqlb.Entity {
	return &mockRecord{}
}

func (m *mockRecords) Append(r sqlb.Entity) {
	*m = append(*m, r.(*mockRecord))
}

type mockDB struct {
	err  error
	data string
	sql  string
}

func (m *mockDB) Query(ctx context.Context, query string, args []interface{}, records sqlb.EntityList) error {
	m.sql = query
	_ = json.Unmarshal([]byte(m.data), records)
	return m.err
}

func (m *mockDB) QueryRow(ctx context.Context, query string, args []interface{}, record sqlb.Entity) error {
	m.sql = query
	_ = json.Unmarshal([]byte(m.data), record)
	return m.err
}

func Test_SelectBuilder_Query(t *testing.T) {
	ctx := context.Background()
	t.Run("no DB provided", func(t *testing.T) {
		err := sqlb.Select("id").FromTable("sample").Query(ctx, &mockRecords{})
		require.EqualError(t, err, "sqlb: no DB was provided to execute the query")
	})

	t.Run("error is propagated properly", func(t *testing.T) {
		db := &mockDB{
			err: errors.New("db error"),
		}

		err := sqlb.MakeSelectBuilder(db, "id").FromTable("sample").Query(ctx, &mockRecords{})
		require.EqualError(t, err, "db error")
	})

	t.Run("query is executed properly", func(t *testing.T) {
		db := &mockDB{
			data: `[{"id":1},{"id":2}]`,
		}
		records := mockRecords{}
		err := sqlb.MakeSelectBuilder(db, "id").FromTable("sample").Query(ctx, &records)
		require.NoError(t, err)
		require.Equal(t, "SELECT id FROM sample", db.sql)
		require.Len(t, records, 2)
		require.Equal(t, 1, records[0].ID)
	})

	t.Run("failed to generate query", func(t *testing.T) {
		records := mockRecords{}
		err := sqlb.Select("id").Where(sqlb.In("id")).Query(ctx, &records)
		require.EqualError(t, err, "values list must not be empty")
	})
}

func Test_SelectBuilder_QueryRow(t *testing.T) {
	ctx := context.Background()
	t.Run("no DB provided", func(t *testing.T) {
		err := sqlb.Select("id").FromTable("sample").QueryRow(ctx, &mockRecord{})
		require.EqualError(t, err, "sqlb: no DB was provided to execute the query")
	})

	t.Run("error is propagated properly", func(t *testing.T) {
		db := &mockDB{
			err: errors.New("db error"),
		}

		err := sqlb.MakeSelectBuilder(db, "id").FromTable("sample").QueryRow(ctx, &mockRecord{})
		require.EqualError(t, err, "db error")
	})

	t.Run("query is executed properly", func(t *testing.T) {
		db := &mockDB{
			data: `{"id":1}`,
		}
		record := mockRecord{}
		err := sqlb.MakeSelectBuilder(db, "id").FromTable("sample").QueryRow(ctx, &record)
		require.NoError(t, err)
		require.Equal(t, "SELECT id FROM sample", db.sql)
		require.Equal(t, 1, record.ID)
	})

	t.Run("failed to generate query", func(t *testing.T) {
		record := mockRecord{}
		err := sqlb.Select("id").Where(sqlb.In("id")).QueryRow(ctx, &record)
		require.EqualError(t, err, "values list must not be empty")
	})
}
