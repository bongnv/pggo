package sqlbuilder

import (
	"io"
)

// ArgAppender is the interface that wraps the Append method.
type ArgAppender interface {
	// Append adds an argument into the existing list.
	Append(values ...interface{})
}

// Condition represents a condition in WHERE clause.
type Condition interface {
	Builder
	conditionOnly()
}

// Equal creates an = condition.
func Equal(column string, value interface{}) Condition {
	return binaryCond{
		operator: "=",
		col:      column,
		value:    placeholder{value: value},
	}
}

// In creates an IN condition.
func In(column string, values ...interface{}) Condition {
	return binaryCond{
		operator: "IN",
		col:      column,
		value:    groupPlaceholder{values: values},
	}
}

// And creates an AND condition.
func And(conds ...Condition) Condition {
	return logicalCond{
		operator: "AND",
		conds:    conds,
	}
}

// Or creates an OR condition.
func Or(conds ...Condition) Condition {
	return logicalCond{
		operator: "OR",
		conds:    conds,
	}
}

type baseCond struct{}

func (c baseCond) conditionOnly() {}

type binaryCond struct {
	baseCond
	operator string
	col      string
	value    Builder
}

func (c binaryCond) Build(sw io.StringWriter, aa ArgAppender) {
	_, _ = sw.WriteString("(")
	_, _ = sw.WriteString(c.col)
	_, _ = sw.WriteString(" ")
	_, _ = sw.WriteString(c.operator)
	_, _ = sw.WriteString(" ")
	c.value.Build(sw, aa)
	_, _ = sw.WriteString(")")
}

type placeholder struct {
	value interface{}
}

func (p placeholder) Build(sw io.StringWriter, aa ArgAppender) {
	_, _ = sw.WriteString("?")
	aa.Append(p.value)
}

type groupPlaceholder struct {
	values []interface{}
}

func (g groupPlaceholder) Build(sw io.StringWriter, aa ArgAppender) {
	count := len(g.values)
	_, _ = sw.WriteString("(")
	for i := 0; i < count; i++ {
		if i > 0 {
			_, _ = sw.WriteString(",")
		}
		_, _ = sw.WriteString("?")
	}
	_, _ = sw.WriteString(")")
	aa.Append(g.values...)
}

type logicalCond struct {
	baseCond
	operator string
	conds    []Condition
}

func (l logicalCond) Build(sw io.StringWriter, aa ArgAppender) {
	if len(l.conds) == 0 {
		return
	}

	if len(l.conds) == 1 {
		l.conds[0].Build(sw, aa)
		return
	}

	_, _ = sw.WriteString("(")
	for i, cond := range l.conds {
		if i > 0 {
			_, _ = sw.WriteString(" ")
			_, _ = sw.WriteString(l.operator)
			_, _ = sw.WriteString(" ")
		}
		cond.Build(sw, aa)
	}
	_, _ = sw.WriteString(")")
}
