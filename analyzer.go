package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strings"
)

func NewAnalyzer(filename string) *Analyzer {
	return &Analyzer{Filename: filename}
}

func (a *Analyzer) Analyze() {
	a.fset = token.NewFileSet()
	var err error
	a.f, err = parser.ParseFile(a.fset, a.Filename, nil, parser.ParseComments)
	CheckError(err)

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
							Analyzer: a,
							Comment:  getComment(genDecl.Doc),
							Name:     typeSpec.Name.Name,
						}
						a.Interfaces = append(a.Interfaces, intf)

						for _, field := range interfaceType.Methods.List {
							m := &Method{
								Analyzer:  a,
								Interface: intf,

								Comment: getComment(field.Doc),
								Name:    field.Names[0].Name,
							}
							intf.Methods = append(intf.Methods, m)

							f, _ := field.Type.(*ast.FuncType)
							for _, in := range f.Params.List {
								param := &VarAndType{
									Analyzer:  a,
									Interface: intf,
									Method:    m,

									Var:  in.Names[0].Name,
									Type: a.Source(in.Type),
								}
								m.Params = append(m.Params, param)
							}
							for i, out := range f.Results.List {
								result := &VarAndType{
									Analyzer:  a,
									Interface: intf,
									Method:    m,

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
		CheckError(err)
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
