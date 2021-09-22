package schema

import "github.com/bongnv/pggo/pkg/sqlb"

// SampleTable defines the schema of sample_table.
var SampleTable = struct {
	sqlb.BaseTable
	ID   string
	Name string
}{
	BaseTable: "sample_table",
	ID:        "id",
	Name:      "name",
}
