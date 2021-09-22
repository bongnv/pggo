package sqlb_test

import (
	"testing"

	"github.com/bongnv/pggo/pkg/sqlb"
)

var globalSQL string

func BenchmarkSelectBuilder_SQL(b *testing.B) {
	b.ReportAllocs()
	var sql string
	for n := 0; n < b.N; n++ {
		sql, _, _ = sqlb.Select("id").FromTable("sample_table").Where(sqlb.Equal("name", "Joe")).SQL()
	}
	globalSQL = sql
}
