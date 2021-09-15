package sqlbuilder_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bongnv/pggo/pkg/sqlbuilder"
)

func Test_SelectBuilder_Build_Happy(t *testing.T) {
	f := sqlbuilder.Factory{}
	personTable := &sqlbuilder.Table{
		Name: "person",
	}
	idCol := &sqlbuilder.Column{
		Table: personTable,
		Name:  "id",
	}
	nameCol := &sqlbuilder.Column{
		Table: personTable,
		Name:  "name",
	}
	sql, args := f.Select(idCol, nameCol).From(personTable).Build()
	require.Equal(t, "SELECT person.id, person.name FROM person", sql)
	require.Len(t, args, 0)
}
