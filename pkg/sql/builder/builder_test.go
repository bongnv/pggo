package builder_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bongnv/pggo/pkg/sql/builder"
	"github.com/bongnv/pggo/pkg/sqlb"
)

type mockRecord struct {
	ID     int
	Name   string
	Number int
}

func (m *mockRecord) GetPointers(cols []string) ([]interface{}, error) {
	if len(cols) == 0 {
		return []interface{}{&m.ID, &m.Name, &m.Number}, nil
	}

	pointers := make([]interface{}, len(cols))
	for i, col := range cols {
		switch col {
		case "id":
			pointers[i] = &m.ID
		case "name":
			pointers[i] = &m.Name
		case "number":
			pointers[i] = &m.Number
		default:
			return nil, fmt.Errorf("%s doesn't exist", col)
		}
	}
	return pointers, nil
}

func (m *mockRecord) GetValues(cols []string) ([]interface{}, error) {
	if len(cols) == 0 {
		return []interface{}{m.ID, m.Name, m.Number}, nil
	}

	values := make([]interface{}, len(cols))
	for i, col := range cols {
		switch col {
		case "id":
			values[i] = m.ID
		case "name":
			values[i] = m.Name
		case "number":
			values[i] = m.Number
		default:
			return nil, fmt.Errorf("%s doesn't exist", col)
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

type mockResult struct {
	err error
}

func (m mockResult) RowsAffected() (int64, error) {
	return 0, m.err
}

func (m mockResult) LastInsertId() (int64, error) {
	return 0, m.err
}

type mockConn struct {
	err    error
	result sql.Result
	sql    string
}

func (m *mockConn) QueryContext(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error) {
	m.sql = sql
	return nil, m.err
}

func (m *mockConn) ExecContext(ctx context.Context, sql string, args ...interface{}) (sql.Result, error) {
	return m.result, m.err
}

func Test_sqlDB_Query(t *testing.T) {
	ctx := context.Background()
	t.Run("with error", func(t *testing.T) {
		m := &mockConn{
			err: errors.New("db error"),
		}
		err := builder.With(m).Select("id").Query(ctx, &mockRecords{})
		require.EqualError(t, err, "db error")
	})
}

func Test_sqlDB_Exec(t *testing.T) {
	ctx := context.Background()
	t.Run("couldn't get affected rows", func(t *testing.T) {
		m := &mockConn{
			result: mockResult{
				err: errors.New("db error"),
			},
		}

		var affectRows int64
		err := builder.With(m).InsertTable("sample_table").Values(1, "Joe").AffectedRows(&affectRows).Exec(ctx)
		require.EqualError(t, err, "db error")
	})
}
