package sqlbuilder_test

import (
	"testing"

	"github.com/bongnv/pggo/pkg/sqlbuilder"
)

var globalSQL string

func BenchmarkSelectBuilder_SQL(b *testing.B) {
	b.ReportAllocs()
	var sql string
	for n := 0; n < b.N; n++ {
		sql, _, _ = sqlbuilder.Select("id").FromTable("sample_table").Where(sqlbuilder.Equal("name", "Joe")).SQL()
	}
	globalSQL = sql
}
