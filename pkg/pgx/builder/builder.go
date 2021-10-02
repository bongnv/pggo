package builder

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgproto3/v2"
	"github.com/jackc/pgx/v4"

	"github.com/bongnv/pggo/pkg/sqlb"
)

// Conn represents a pgx DB connection.
type Conn interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

// With creates a builder factory from a DB connection.
func With(conn Conn) sqlb.Factory {
	return sqlb.Factory{
		DB: pgxDB{
			conn: conn,
		},
	}
}

type pgxDB struct {
	conn Conn
}

func (db pgxDB) Query(ctx context.Context, query string, args []interface{}, records sqlb.EntityList) error {
	rows, err := db.conn.Query(ctx, query, args...)
	if err != nil {
		return err
	}

	fields := rows.FieldDescriptions()
	defer rows.Close()

	for rows.Next() {
		item := records.New()
		pointers, err := buildPointers(item, fields)
		if err != nil {
			return err
		}

		if err := rows.Scan(pointers...); err != nil {
			return err
		}

		records.Append(item)
	}

	return rows.Err()
}

func (db pgxDB) QueryRow(ctx context.Context, query string, args []interface{}, record sqlb.Entity) error {
	rows, err := db.conn.Query(ctx, query, args...)
	if err != nil {
		return err
	}

	if !rows.Next() {
		if rows.Err() == nil {
			return pgx.ErrNoRows
		}

		return rows.Err()
	}

	defer rows.Close()

	fields := rows.FieldDescriptions()
	pointers, err := buildPointers(record, fields)
	if err != nil {
		return err
	}

	if err := rows.Scan(pointers...); err != nil {
		return err
	}

	return rows.Err()
}

func (db pgxDB) Exec(ctx context.Context, sql string, args []interface{}, affectedRows *int64) error {
	res, err := db.conn.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	if affectedRows != nil {
		*affectedRows = res.RowsAffected()
	}

	return nil
}
func buildPointers(record sqlb.Entity, fields []pgproto3.FieldDescription) ([]interface{}, error) {
	cols := make([]string, len(fields))

	for i, f := range fields {
		cols[i] = string(f.Name)
	}

	return record.GetPointers(cols)
}
