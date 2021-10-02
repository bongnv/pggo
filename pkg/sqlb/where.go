package sqlb

import "io"

type whereClause struct {
	cond Condition
}

func (w whereClause) Build(sb io.StringWriter, aa Placeholders) error {
	_, _ = sb.WriteString(" WHERE ")
	return w.cond.Build(sb, aa)
}
