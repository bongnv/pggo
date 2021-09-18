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
		expectedArgs  sqlbuilder.ArgumentList
	}{
		"equal": {
			createCond: func() sqlbuilder.Condition {
				return sqlbuilder.Equal("id", 10)
			},
			expectedQuery: "(id = $1)",
			expectedArgs:  []interface{}{10},
		},
		"in": {
			createCond: func() sqlbuilder.Condition {
				return sqlbuilder.In("id", 1, 2, 3, 4)
			},
			expectedQuery: "(id IN ($1,$2,$3,$4))",
			expectedArgs:  []interface{}{1, 2, 3, 4},
		},
		"and": {
			createCond: func() sqlbuilder.Condition {
				cond1 := sqlbuilder.Equal("id", 10)
				cond2 := sqlbuilder.Equal("name", "Joe")
				return sqlbuilder.And(cond1, cond2)
			},
			expectedQuery: "((id = $1) AND (name = $2))",
			expectedArgs:  []interface{}{10, "Joe"},
		},
		"or": {
			createCond: func() sqlbuilder.Condition {
				cond1 := sqlbuilder.Equal("id", 10)
				cond2 := sqlbuilder.Equal("name", "Joe")
				return sqlbuilder.Or(cond1, cond2)
			},
			expectedQuery: "((id = $1) OR (name = $2))",
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
			args := sqlbuilder.ArgumentList{}
			tc.createCond().Build(sb, &args)
			require.Equal(t, tc.expectedQuery, sb.String())
			if tc.expectedArgs == nil {
				require.Empty(t, args)
			} else {
				require.Equal(t, tc.expectedArgs, args)
			}
		})
	}
}
