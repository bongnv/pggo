package template_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bongnv/pggo/internal/template"
)

type mockData struct {
	Name string
}

func TestExecute(t *testing.T) {
	var buf bytes.Buffer
	err := template.Execute(&buf, "unit_test.tmpl", &mockData{Name: "TestExecute"})
	require.NoError(t, err)
	require.Equal(t, "This is executed with data: TestExecute.\n", buf.String())
}
