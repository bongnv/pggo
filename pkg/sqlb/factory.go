package sqlb

import (
	"context"
	"errors"
)

// DB represents a DB for running queries one.
type DB interface {
	Query(ctx context.Context, query string, args []interface{}, records EntityList) error
	QueryRow(ctx context.Context, query string, args []interface{}, record Entity) error
	Exec(ctx context.Context, sql string, arg []interface{}, affectedRows *int64) error
}

// Factory is a builder factory for creating builders.
type Factory struct {
	DB DB
}

// Select starts a new SELECT query.
func (f Factory) Select(cols ...string) *SelectBuilder {
	return &SelectBuilder{
		cols: cols,
		db:   f.DB,
	}
}

// Insert starts a new INSERT query.
func (f Factory) Insert(table Table) *InsertBuilder {
	return &InsertBuilder{
		table: table,
		db:    f.DB,
	}
}

// Insert starts a new INSERT query with a table name.
func (f Factory) InsertTable(tableName string) *InsertBuilder {
	return f.Insert(BaseTable(tableName))
}

// DefaultFactory is the default factory.
var DefaultFactory = Factory{
	DB: noopDB{},
}

type noopDB struct{}

func (noopDB) Query(ctx context.Context, query string, args []interface{}, records EntityList) error {
	return errors.New("sqlb: no DB was provided to execute the query")
}

func (noopDB) QueryRow(ctx context.Context, query string, args []interface{}, record Entity) error {
	return errors.New("sqlb: no DB was provided to execute the query")
}

func (noopDB) Exec(ctx context.Context, sql string, arg []interface{}, affectedRows *int64) error {
	return errors.New("sqlb: no DB was provided to execute the query")
}
