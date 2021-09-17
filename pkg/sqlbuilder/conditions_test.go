package sqlbuilder_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bongnv/pggo/pkg/sqlbuilder"
)

func Test_Conditions(t *testing.T) {
	cases := map[string]struct {
		createCond    func() sqlbuilder.Condition
		expectedQuery string
		expectedArgs  []interface{}
	}{
		"equal": {
			createCond: func() sqlbuilder.Condition {
				return sqlbuilder.Equal("id", 10)
			},
			expectedQuery: "(id = ?)",
			expectedArgs:  []interface{}{10},
		},
		"in": {
			createCond: func() sqlbuilder.Condition {
				return sqlbuilder.In("id", 1, 2, 3, 4)
			},
			expectedQuery: "(id IN (?,?,?,?))",
			expectedArgs:  []interface{}{1, 2, 3, 4},
		},
		"and": {
			createCond: func() sqlbuilder.Condition {
				cond1 := sqlbuilder.Equal("id", 10)
				cond2 := sqlbuilder.Equal("name", "Joe")
				return sqlbuilder.And(cond1, cond2)
			},
			expectedQuery: "((id = ?) AND (name = ?))",
			expectedArgs:  []interface{}{10, "Joe"},
		},
		"or": {
			createCond: func() sqlbuilder.Condition {
				cond1 := sqlbuilder.Equal("id", 10)
				cond2 := sqlbuilder.Equal("name", "Joe")
				return sqlbuilder.Or(cond1, cond2)
			},
			expectedQuery: "((id = ?) OR (name = ?))",
			expectedArgs:  []interface{}{10, "Joe"},
		},
		"empty-and": {
			createCond: func() sqlbuilder.Condition {
				return sqlbuilder.And()
			},
			expectedQuery: "",
			expectedArgs:  nil,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			sb := &strings.Builder{}
			args := &sqlbuilder.ArgumentList{}
			tc.createCond().Build(sb, args)
			require.Equal(t, tc.expectedQuery, sb.String())
			require.Equal(t, tc.expectedArgs, args.Args)
		})
	}
}
