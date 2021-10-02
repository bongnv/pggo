package builder_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"

	"github.com/bongnv/pggo/pkg/pgx/builder"
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

type mockConn struct {
	err        error
	sql        string
	commandTag pgconn.CommandTag
}

func (m *mockConn) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	m.sql = sql
	return nil, m.err
}

func (m *mockConn) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return m.commandTag, m.err
}

func Test_pgxDB_Query(t *testing.T) {
	ctx := context.Background()
	t.Run("with error", func(t *testing.T) {
		m := &mockConn{
			err: errors.New("db error"),
		}
		err := builder.With(m).Select("id").Query(ctx, &mockRecords{})
		require.EqualError(t, err, "db error")
	})
}
