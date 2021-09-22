package builder

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bongnv/pggo/pkg/sqlb"
)

// Conn represents a DB connection.
type Conn interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
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

func (db sqlDB) Query(ctx context.Context, query string, args []interface{}, records sqlb.Recordables) error {
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

func (db sqlDB) QueryRow(ctx context.Context, query string, args []interface{}, record sqlb.Recordable) error {
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

	pointers, err := buildPointers(record, fields)
	if err != nil {
		return err
	}

	if err := rows.Scan(pointers...); err != nil {
		return err
	}

	return rows.Err()
}

func buildPointers(record sqlb.Recordable, fields []string) ([]interface{}, error) {
	pointers := make([]interface{}, 0, len(fields))
	for _, field := range fields {
		pointer := record.GetPointer(field)
		if pointer == nil {
			return nil, fmt.Errorf("builder: %s is not found", field)
		}
		pointers = append(pointers, pointer)
	}

	return pointers, nil
}
