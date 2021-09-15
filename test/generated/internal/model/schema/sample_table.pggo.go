package schema

import "github.com/bongnv/pggo/pkg/sqlbuilder"

// SampleTable defines the schema of sample_table.
var SampleTable = struct {
	sqlbuilder.Table
	ID   sqlbuilder.Column
	Name sqlbuilder.Column
}{
	Table: sqlbuilder.Table{
		Name: "sample_table",
	},
}

func init() {
	SampleTable.ID = sqlbuilder.Column{
		Table: SampleTable,
		Name:  "id",
	}

	SampleTable.Name = sqlbuilder.Column{
		Table: SampleTable,
		Name:  "name",
	}
}
