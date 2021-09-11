package loader

import (
	"context"

	"github.com/jackc/pgx/v4"

	"github.com/bongnv/pggo/internal/generator"
)

// PostgreSQLLoader is an implementation to load schema from PostgreSQL DB.
type PostgreSQLLoader struct {
	URL string
}

func (l PostgreSQLLoader) Load() (generator.Schema, error) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, l.URL)
	if err != nil {
		return generator.Schema{}, err
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		return generator.Schema{}, err
	}

	defer rows.Close()

	tables := []generator.Table{}

	for rows.Next() {
		table := generator.Table{}
		if err := rows.Scan(&table.Name); err != nil {
			return generator.Schema{}, err
		}

		tables = append(tables, table)
	}

	return generator.Schema{
		Tables: tables,
	}, nil
}
