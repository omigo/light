package goparser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os/exec"
	"strings"

	"github.com/arstd/log"
)

func Parse(src string) *Store {

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, src, nil, parser.ParseComments)
	if err != nil {
		log.Panic(err)
	}
	// ast.Print(fset, f)

	store := &Store{
		Source:  src,
		Package: f.Name.Name,
		Imports: map[string]string{
			"bytes": "",
			"github.com/arstd/light/light": "",
			"github.com/arstd/log":         "",
		},
	}

	goBuild(src)

	extractDocs(store, f)
	parseTypes(store, fset, f)

	return store
}

func goBuild(src string) {
	// log.Debugf("go build -i -v  %s", goFile)
	cmd := exec.Command("go", "build", "-i", "-v", src)
	out, err := cmd.CombinedOutput()
	if bytes.HasSuffix(out, []byte("command-line-arguments\n")) {
		fmt.Printf("%s", out[:len(out)-23])
	} else {
		fmt.Printf("%s", out)
	}
	if err != nil {
		log.Panic(err)
	}
}

func extractDocs(store *Store, f *ast.File) {
	for _, decl := range f.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
						if store.Name != "" {
							panic("one file must contains one interface only")
						}

						store.Name = typeSpec.Name.Name
						for _, field := range interfaceType.Methods.List {
							m := &Method{
								Store: store,
								Name:  field.Names[0].Name,
								Doc:   getDoc(field.Doc),
							}
							store.Methods = append(store.Methods, m)
						}
					}
				}
			}
		}
	}
}

func getDoc(cg *ast.CommentGroup) (comment string) {
	if cg == nil {
		return ""
	}
	for _, c := range cg.List {
		comment += strings.TrimSpace(c.Text[2:]) + " " // remove `//`
	}
	return strings.TrimSpace(comment)
}

func parseTypes(store *Store, fset *token.FileSet, f *ast.File) {
	info := types.Info{Defs: make(map[*ast.Ident]types.Object)}
	conf := types.Config{Importer: importer.Default()}
	log.Fataln(conf.Check(store.Name, fset, []*ast.File{f}, &info))

	for k, obj := range info.Defs {
		if k.Obj != nil && k.Name == store.Name && k.Obj.Kind == ast.Typ {
			// get method name and params/returns
			if itfType, ok := obj.Type().Underlying().(*types.Interface); ok {
				for i := 0; i < itfType.NumMethods(); i++ {
					x := itfType.Method(i)
					m := store.MethodByName(x.Name())
					y := x.Type().(*types.Signature)
					m.Params = &Tuple{store, y.Params()}
					m.Results = &Tuple{store, y.Results()}
				}
			}
		}
	}
}
