package builder

import (
	"context"
	"database/sql"

	"github.com/bongnv/pggo/pkg/sqlb"
)

// Conn represents a DB connection.
type Conn interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, sql string, args ...interface{}) (sql.Result, error)
}

// With creates a builder factory to start building queries.
func With(conn Conn) sqlb.Factory {
	return sqlb.Factory{
		DB: &sqlDB{
			conn: conn,
		},
	}
}

type sqlDB struct {
	conn Conn
}

func (db sqlDB) Query(ctx context.Context, query string, args []interface{}, records sqlb.EntityList) error {
	rows, err := db.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}

	fields, err := rows.Columns()
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		item := records.New()
		pointers, err := item.GetPointers(fields)
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

func (db sqlDB) QueryRow(ctx context.Context, query string, args []interface{}, record sqlb.Entity) error {
	rows, err := db.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}

	if !rows.Next() {
		if rows.Err() == nil {
			return sql.ErrNoRows
		}

		return rows.Err()
	}

	defer rows.Close()

	fields, err := rows.Columns()
	if err != nil {
		return err
	}

	pointers, err := record.GetPointers(fields)
	if err != nil {
		return err
	}

	if err := rows.Scan(pointers...); err != nil {
		return err
	}

	return rows.Err()
}

func (db sqlDB) Exec(ctx context.Context, sql string, args []interface{}, affectedRows *int64) error {
	res, err := db.conn.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	if affectedRows != nil {
		if *affectedRows, err = res.RowsAffected(); err != nil {
			return err
		}
	}

	return nil
}
