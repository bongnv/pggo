package sqlb

import (
	"errors"
	"io"
)

// Placeholders is the interface that wraps the Append method.
type Placeholders interface {
	// Append adds an argument into the existing and generate a placeholder for it..
	Append(values ...interface{}) string
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

func (c binaryCond) Build(sw io.StringWriter, aa Placeholders) error {
	_, _ = sw.WriteString("(")
	_, _ = sw.WriteString(c.col)
	_, _ = sw.WriteString(" ")
	_, _ = sw.WriteString(c.operator)
	_, _ = sw.WriteString(" ")
	if err := c.value.Build(sw, aa); err != nil {
		return err
	}
	_, _ = sw.WriteString(")")
	return nil
}

type placeholder struct {
	value interface{}
}

func (p placeholder) Build(sw io.StringWriter, aa Placeholders) error {
	_, _ = sw.WriteString(aa.Append(p.value))
	return nil
}

type groupPlaceholder struct {
	values []interface{}
}

func (g groupPlaceholder) Build(sw io.StringWriter, aa Placeholders) error {
	if len(g.values) == 0 {
		return errors.New("values list must not be empty")
	}

	_, _ = sw.WriteString("(")
	for i, v := range g.values {
		if i > 0 {
			_, _ = sw.WriteString(",")
		}
		_, _ = sw.WriteString(aa.Append(v))
	}
	_, _ = sw.WriteString(")")

	return nil
}

type logicalCond struct {
	baseCond
	operator string
	conds    []Condition
}

func (l logicalCond) Build(sw io.StringWriter, aa Placeholders) error {
	if len(l.conds) == 0 {
		return errors.New("conditions list must not be empty")
	}

	if len(l.conds) == 1 {
		return l.conds[0].Build(sw, aa)
	}

	_, _ = sw.WriteString("(")
	for i, cond := range l.conds {
		if i > 0 {
			_, _ = sw.WriteString(" ")
			_, _ = sw.WriteString(l.operator)
			_, _ = sw.WriteString(" ")
		}
		if err := cond.Build(sw, aa); err != nil {
			return err
		}
	}
	_, _ = sw.WriteString(")")
	return nil
}
