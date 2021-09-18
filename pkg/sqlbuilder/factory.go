package sqlbuilder

import (
	"context"
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
func (f Factory) Select(cols ...string) SelectBuilder {
	return SelectBuilder{
		cols: cols,
		db:   f.DB,
	}
}
