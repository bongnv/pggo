package sqlbuilder

import "io"

// ArgAppender is the interface that wraps the Append method.
type ArgAppender interface {
	// Append adds an argument into the existing list.
	Append(values ...interface{})
}

// Condition represents a condition in WHERE clause.
type Condition func(sb io.StringWriter, aa ArgAppender)

// Equal creates an = condition.
func Equal(column string, value interface{}) Condition {
	return func(sb io.StringWriter, aa ArgAppender) {
		_, _ = sb.WriteString("(")
		_, _ = sb.WriteString(column)
		_, _ = sb.WriteString(" = ?)")
		aa.Append(value)
	}
}

// In creates an IN condition.
func In(column string, values ...interface{}) Condition {
	return func(sb io.StringWriter, aa ArgAppender) {
		_, _ = sb.WriteString("(")
		_, _ = sb.WriteString(column)
		_, _ = sb.WriteString(" IN (")
		for i := range values {
			if i > 0 {
				_, _ = sb.WriteString(",")
			}
			_, _ = sb.WriteString("?")
		}
		_, _ = sb.WriteString(")")
		aa.Append(values...)
	}
}
