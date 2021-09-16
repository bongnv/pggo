package sqlbuilder

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Equal(t *testing.T) {
	sb := &strings.Builder{}
	args := &argumentList{}
	Equal("id", 10)(sb, args)
	require.Equal(t, "(id = ?)", sb.String())
	require.Equal(t, []interface{}{10}, args.args)
}

func Test_In(t *testing.T) {
	sb := &strings.Builder{}
	args := &argumentList{}
	In("id", 1, 2, 3, 4)(sb, args)
	require.Equal(t, "(id IN (?,?,?,?)", sb.String())
	require.Equal(t, []interface{}{1, 2, 3, 4}, args.args)
}
