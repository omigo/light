package main

import (
	"bytes"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/arstd/log"
)

func NewAnalyzer(filename string) *Analyzer {
	return &Analyzer{Filename: filename}
}

func (a *Analyzer) Analyze() {
	a.fset = token.NewFileSet()
	var err error
	a.f, err = parser.ParseFile(a.fset, a.Filename, nil, parser.ParseComments)
	checkError(err)

	// ast.Print(a.fset, a.f)
	// format.Node(os.Stdout, fset, f)

	a.Package = a.f.Name.Name // package name

	for _, imp := range a.f.Imports {
		a.Imports = append(a.Imports, a.Source(imp))
	}

	for _, decl := range a.f.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			if genDecl.Tok == token.TYPE && len(genDecl.Specs) == 1 {
				if typeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec); ok {
					if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
						intf := &Interface{
							a:       a,
							Comment: getComment(genDecl.Doc),
							Name:    typeSpec.Name.Name,
						}
						a.Interfaces = append(a.Interfaces, intf)

						for _, field := range interfaceType.Methods.List {
							m := &Method{
								a: a,
								i: intf,

								Comment: getComment(field.Doc),
								Name:    field.Names[0].Name,
							}
							intf.Methods = append(intf.Methods, m)

							f, _ := field.Type.(*ast.FuncType)
							for _, in := range f.Params.List {
								param := &VarAndType{
									a: a,
									i: intf,
									m: m,

									Var:  in.Names[0].Name,
									Type: a.Source(in.Type),
								}
								m.Params = append(m.Params, param)
								deepParse(in.Type, param)
							}
							for i, out := range f.Results.List {
								result := &VarAndType{
									a: a,
									i: intf,
									m: m,

									// Var:  out.Names[0].Name,
									Type: a.Source(out.Type),
								}
								if len(out.Names) == 0 {
									result.Var = genVarByType(result.Type)
								}

								if result.Type == "error" && (i != f.Results.NumFields()-1) {
									panic(" for method '" + m.Name + "', 'error' must be last return value")
								}

								m.Results = append(m.Results, result)

								deepParse(out.Type, result)
							}
						}
					}
				}
			}
		}
	}
}

func (a *Analyzer) Source(node ast.Node) string {
	var buf bytes.Buffer
	err := format.Node(&buf, a.fset, node)
	if err != nil {
		checkError(err)
	}
	return buf.String()
}

func getComment(cg *ast.CommentGroup) (comment string) {
	if cg == nil {
		return ""
	}
	for _, c := range cg.List {
		comment += strings.TrimSpace(c.Text[2:]) + " " // remove `//`
	}
	return strings.TrimSpace(comment)
}

func genVarByType(t string) string {
	if t == "error" {
		return "err"
	}
	t = strings.TrimLeft(t, "[]*")
	t = t[strings.LastIndex(t, ".")+1:]
	return "v" + strings.ToLower(t[0:1])
}

func deepParse(expr ast.Expr, vt *VarAndType) string {
	switch e := expr.(type) {
	case *ast.ArrayType:
		vt.Slice = "[]"
		return "[]" + deepParse(e.Elt, vt)

	case *ast.StarExpr:
		vt.Star = "*"
		return "*" + deepParse(e.X, vt)

	case *ast.Ident:
		return e.Name

	case *ast.SelectorExpr:
		x := deepParse(e.X, vt)
		vt.Pkg = x
		parseSelector(vt.a.GetPath(x), e.Sel.Name, vt)
		return x + "." + e.Sel.Name

	default:
		log.Warnf("unimplemented %T", e)
		return ""
	}
}

func parseSelector(path, sel string, vt *VarAndType) {
	if path == "database/sql" && sel == "Tx" {
		return
	}

	pkg, err := build.Import(path, "", 0)
	checkError(err)

	fset := token.NewFileSet()
	for _, file := range pkg.GoFiles {
		f, err := parser.ParseFile(fset, filepath.Join(pkg.Dir, file), nil, 0)
		checkError(err)

		ast.Print(fset, f)

		for _, decl := range f.Decls {
			decl, ok := decl.(*ast.GenDecl)
			if !ok || decl.Tok != token.TYPE {
				continue
			}
			for _, spec := range decl.Specs {
				switch t := spec.(type) {
				case *ast.TypeSpec:
					switch u := t.Type.(type) {
					case *ast.StructType:
						vt.Fields = []*VarAndType{}
						for _, f := range u.Fields.List {
							fvt := &VarAndType{}
							vt.Fields = append(vt.Fields, fvt)

							fvt.Var = deepParse(f.Names[0], fvt)
							// fvt.Type = deepParse(f.Type, f)
						}
					}
				}
			}
		}

	}

}
