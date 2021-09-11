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

// Load connects to the given URL to load DB schema.
func (l PostgreSQLLoader) Load() (*generator.Schema, error) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, l.URL)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tables := map[string]*generator.Table{}

	for rows.Next() {
		table := &generator.Table{}
		if err := rows.Scan(&table.Name); err != nil {
			return nil, err
		}

		tables[table.Name] = table
	}

	if err := fetchColumns(conn, tables); err != nil {
		return nil, err
	}

	return &generator.Schema{
		Tables: getTables(tables),
	}, nil
}

func getTables(tables map[string]*generator.Table) []*generator.Table {
	results := make([]*generator.Table, 0, len(tables))
	for _, v := range tables {
		results = append(results, v)
	}

	return results
}

func fetchColumns(conn *pgx.Conn, tables map[string]*generator.Table) error {
	ctx := context.Background()
	rows, err := conn.Query(ctx, "SELECT table_name, column_name, is_nullable, data_type FROM information_schema.columns WHERE table_schema = 'public'")
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		column := &generator.Column{}
		var nullable string
		var tableName string
		if err := rows.Scan(&tableName, &column.Name, &nullable, &column.DataType); err != nil {
			return err
		}

		column.Nullable = nullable == "YES"
		tables[tableName].Columns = append(tables[tableName].Columns, column)
	}

	return nil
}
