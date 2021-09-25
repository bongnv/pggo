package builder_test

import (
	"context"
	"errors"
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

func (m *mockRecord) GetPointer(col string) interface{} {
	switch col {
	case "id":
		return &m.ID
	case "name":
		return &m.Name
	case "number":
		return &m.Number
	default:
		return nil
	}
}

type mockRecords []*mockRecord

func (m mockRecords) New() sqlb.Recordable {
	return &mockRecord{}
}

func (m *mockRecords) Append(r sqlb.Recordable) {
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
