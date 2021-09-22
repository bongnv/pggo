package sqlb_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bongnv/pggo/pkg/sqlb"
)

func Test_Conditions(t *testing.T) {
	cases := map[string]struct {
		createCond    func() sqlb.Condition
		expectedQuery string
		expectedArgs  sqlb.ArgumentList
	}{
		"equal": {
			createCond: func() sqlb.Condition {
				return sqlb.Equal("id", 10)
			},
			expectedQuery: "(id = $1)",
			expectedArgs:  []interface{}{10},
		},
		"in": {
			createCond: func() sqlb.Condition {
				return sqlb.In("id", 1, 2, 3, 4)
			},
			expectedQuery: "(id IN ($1,$2,$3,$4))",
			expectedArgs:  []interface{}{1, 2, 3, 4},
		},
		"and": {
			createCond: func() sqlb.Condition {
				cond1 := sqlb.Equal("id", 10)
				cond2 := sqlb.Equal("name", "Joe")
				return sqlb.And(cond1, cond2)
			},
			expectedQuery: "((id = $1) AND (name = $2))",
			expectedArgs:  []interface{}{10, "Joe"},
		},
		"or": {
			createCond: func() sqlb.Condition {
				cond1 := sqlb.Equal("id", 10)
				cond2 := sqlb.Equal("name", "Joe")
				return sqlb.Or(cond1, cond2)
			},
			expectedQuery: "((id = $1) OR (name = $2))",
			expectedArgs:  []interface{}{10, "Joe"},
		},
		"empty-and": {
			createCond: func() sqlb.Condition {
				return sqlb.And()
			},
			expectedQuery: "",
			expectedArgs:  nil,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			sb := &strings.Builder{}
			args := sqlb.ArgumentList{}
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
