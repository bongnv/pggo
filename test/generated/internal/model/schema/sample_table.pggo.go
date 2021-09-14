package schema

import "github.com/bongnv/pggo/pkg/sqlschema"

// SampleTable defines the schema of sample_table.
var SampleTable = struct {
	sqlschema.Table
	ID   sqlschema.Column
	Name sqlschema.Column
}{
	Table: sqlschema.Table{
		Name: "sample_table",
	},
}

func init() {
	SampleTable.ID = sqlschema.Column{
		Table: SampleTable,
		Name:  "id",
	}

	SampleTable.Name = sqlschema.Column{
		Table: SampleTable,
		Name:  "name",
	}
}
