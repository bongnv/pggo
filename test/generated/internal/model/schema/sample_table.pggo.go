package schema

import "github.com/bongnv/pggo/pkg/sqlbuilder"

// SampleTable defines the schema of sample_table.
var SampleTable = struct {
	sqlbuilder.BaseTable
	ID   string
	Name string
}{
	BaseTable: "sample_table",
	ID:        "id",
	Name:      "name",
}
