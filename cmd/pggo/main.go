package main

import (
	"github.com/alecthomas/kong"

	"github.com/bongnv/pggo/internal/generator"
	"github.com/bongnv/pggo/internal/loader"
	"github.com/bongnv/pggo/internal/writer"
)

var cli struct {
	Dir   string `kong:"optional,name='dir',short='d',default='.',help='Directory for output files'"`
	URL   string `kong:"required,name='url',short='u',help='Connection URL to PostgreSQL server'"`
	Table string `kong:"optional,name='table',short='t',help='Name of the table for generating code'"`
}

func main() {
	ctx := kong.Parse(
		&cli,
		kong.Name("gogen"),
		kong.Description("A code generation tool using Go template"),
	)

	loader := loader.PostgreSQLLoader{
		URL: cli.URL,
	}

	writer := writer.FileWriter{
		Dir: cli.Dir,
	}

	gen := generator.Generator{
		SchemaLoader: loader,
		Table:        cli.Table,
		Writer:       writer,
	}

	ctx.FatalIfErrorf(gen.Generate())
}
