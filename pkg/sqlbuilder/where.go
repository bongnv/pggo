package sqlbuilder

import "io"

type whereClause struct {
	cond Condition
}

func (w whereClause) Build(sb io.StringWriter, aa Placeholders) {
	_, _ = sb.WriteString(" WHERE ")
	w.cond.Build(sb, aa)
}
