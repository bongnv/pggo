package schema

import "github.com/bongnv/pggo/pkg/sqlbuilder"

// SampleTable defines the schema of sample_table.
var SampleTable = struct {
	sqlbuilder.Table
	ID   sqlbuilder.Column
	Name sqlbuilder.Column
}{
	Table: sqlbuilder.BaseTable{
		Name: "sample_table",
	},
}

func init() {
	SampleTable.ID = sqlbuilder.BaseColumn{
		Table: SampleTable,
		Name:  "id",
	}

	SampleTable.Name = sqlbuilder.BaseColumn{
		Table: SampleTable,
		Name:  "name",
	}
}
