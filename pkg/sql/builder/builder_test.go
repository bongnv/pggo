package builder_test

import (
	"context"
	"database/sql"
	"errors"
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
	err error
	sql string
}

func (m *mockConn) QueryContext(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error) {
	m.sql = sql
	return nil, m.err
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
