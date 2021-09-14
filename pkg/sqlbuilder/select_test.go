package sqlbuilder_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bongnv/pggo/pkg/sqlbuilder"
)

type mockTable struct {
	name string
}

func (t mockTable) SQL() string {
	return t.name
}

type mockColumn struct {
	table string
	name  string
}

func (c mockColumn) SQL() string {
	return c.table + "." + c.name
}

func (c mockColumn) TableName() string {
	return c.table
}

func Test_SelectBuilder_Build_Happy(t *testing.T) {
	f := sqlbuilder.Factory{}
	personTable := &mockTable{
		name: "person",
	}
	idCol := &mockColumn{
		name:  "id",
		table: "person",
	}
	nameCol := &mockColumn{
		name:  "name",
		table: "person",
	}
	sql, args := f.Select(idCol, nameCol).From(personTable).Build()
	require.Equal(t, "SELECT person.id, person.name FROM person", sql)
	require.Len(t, args, 0)
}
