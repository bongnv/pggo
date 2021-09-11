package main

import (
	"github.com/alecthomas/kong"

	"github.com/bongnv/pggo/internal/generator"
	"github.com/bongnv/pggo/internal/loader"
)

var cli struct {
	URL string `kong:"required,name='url',short='u',help='Connection URL to PostgreSQL server'"`
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

	gen := generator.Generator{
		SchemaLoader: loader,
	}

	ctx.FatalIfErrorf(gen.Generate())
}
