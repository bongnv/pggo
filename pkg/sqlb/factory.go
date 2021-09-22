package sqlb

import (
	"context"
	"errors"
)

// DB represents a DB for running queries one.
type DB interface {
	Query(ctx context.Context, query string, args []interface{}, records Recordables) error
	QueryRow(ctx context.Context, query string, args []interface{}, record Recordable) error
}

// Factory is a builder factory for creating builders.
type Factory struct {
	DB DB
}

// Select starts a new SELECT query.
func (f Factory) Select(cols ...string) *SelectBuilder {
	db := f.DB
	if db == nil {
		db = noopDB{}
	}

	return &SelectBuilder{
		cols: cols,
		db:   db,
	}
}

type noopDB struct{}

func (noopDB) Query(ctx context.Context, query string, args []interface{}, records Recordables) error {
	return errors.New("sqlb: no DB was provided to execute the query")
}

func (noopDB) QueryRow(ctx context.Context, query string, args []interface{}, record Recordable) error {
	return errors.New("sqlb: no DB was provided to execute the query")
}
