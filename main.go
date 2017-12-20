package main

import (
	"go/parser"
	"go/token"
	"os"

	"github.com/arstd/log"
)

func main() {
	src := os.Getenv("GOFILE")
	if src == "" {
		src = "/Users/Arstd/Reposits/src/github.com/arstd/light/example/store/user.go"
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, src, nil, parser.ParseComments)
	if err != nil {
		log.Panic(err)
	}
	// ast.Print(fset, f)

	store := &Store{Package: f.Name.Name, Imports: map[string]string{}}

	goBuild(src)

	extractDocs(store, f)
	parseTypes(store, fset, f)

	// log.JSONIndent(store)

	Generate(store)
}
