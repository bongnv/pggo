package template

import (
	"embed"
	"io"
	"text/template"
)

//go:embed *.tmpl
var tmplFiles embed.FS
var rootTemplate = getRootTemplate()

// Execute executes a template given a name.
func Execute(out io.Writer, name string, data interface{}) error {
	return rootTemplate.ExecuteTemplate(out, name, data)
}

func getRootTemplate() *template.Template {
	rootTemplate, err := template.ParseFS(tmplFiles, "*")
	if err != nil {
		panic(err)
	}

	return rootTemplate
}
