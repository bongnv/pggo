package sqlbuilder_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bongnv/pggo/pkg/sqlbuilder"
)

func Test_Equal(t *testing.T) {
	sb := &strings.Builder{}
	args := &sqlbuilder.ArgumentList{}
	sqlbuilder.Equal("id", 10).Build(sb, args)
	require.Equal(t, "(id = ?)", sb.String())
	require.Equal(t, []interface{}{10}, args.Args)
}

func Test_In(t *testing.T) {
	sb := &strings.Builder{}
	args := &sqlbuilder.ArgumentList{}
	sqlbuilder.In("id", 1, 2, 3, 4).Build(sb, args)
	require.Equal(t, "(id IN (?,?,?,?))", sb.String())
	require.Equal(t, []interface{}{1, 2, 3, 4}, args.Args)
}

func Test_And(t *testing.T) {
	sb := &strings.Builder{}
	args := &sqlbuilder.ArgumentList{}

	cond1 := sqlbuilder.Equal("id", 10)
	cond2 := sqlbuilder.Equal("name", "Joe")
	sqlbuilder.And(cond1, cond2).Build(sb, args)
	require.Equal(t, "((id = ?) AND (name = ?))", sb.String())
	require.Equal(t, []interface{}{10, "Joe"}, args.Args)
}

func Test_Or(t *testing.T) {
	sb := &strings.Builder{}
	args := &sqlbuilder.ArgumentList{}

	cond1 := sqlbuilder.Equal("id", 10)
	cond2 := sqlbuilder.Equal("name", "Joe")
	sqlbuilder.Or(cond1, cond2).Build(sb, args)
	require.Equal(t, "((id = ?) OR (name = ?))", sb.String())
	require.Equal(t, []interface{}{10, "Joe"}, args.Args)
}
