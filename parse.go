package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os/exec"
	"regexp"
	"strings"

	"github.com/arstd/log"
)

func parseGoFile(pkg *Package) {
	defer func() { parsed = nil }()

	goBuild(pkg.Source)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, pkg.Source, nil, parser.ParseComments)
	if err != nil {
		log.Panic(err)
	}

	pkg.Name = f.Name.Name
	parseImports(pkg, f.Imports)

	parseComments(pkg, f)

	// ast.Print(fset, f)
	// printer.Fprint(os.Stdout, fset, f)

	parseTypes(pkg, fset, f)
}

func goBuild(goFile string) {
	log.Debugf("go build -i -v  %s", goFile)
	cmd := exec.Command("go", "build", "-i", "-v", goFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Panic(err)
	}
	if bytes.HasSuffix(out, []byte("command-line-arguments\n")) {
		fmt.Printf("%s", out[:len(out)-23])
	} else {
		fmt.Printf("%s", out)
	}
}

func parseComments(pkg *Package, f *ast.File) {
	for _, decl := range f.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
						itf := &Interface{
							Doc:  getDoc(genDecl.Doc),
							Name: typeSpec.Name.Name,
						}
						pkg.Interfaces = append(pkg.Interfaces, itf)

						for _, field := range interfaceType.Methods.List {
							m := &Method{
								Name: field.Names[0].Name,
								Doc:  getDoc(field.Doc),
							}
							itf.Methods = append(itf.Methods, m)
						}
					}
				}
			}
		}
	}
}

func parseTypes(pkg *Package, fset *token.FileSet, f *ast.File) {
	info := types.Info{
		Defs: make(map[*ast.Ident]types.Object),
	}
	conf := types.Config{Importer: importer.Default()}
	_, err := conf.Check(pkg.Name, fset, []*ast.File{f}, &info)
	if err != nil {
		panic(err)
	}

	for k, obj := range info.Defs {
		if k.Obj == nil || k.Obj.Kind != ast.Typ {
			continue
		}
		var itf *Interface
		for _, x := range pkg.Interfaces {
			if x.Name == k.Name {
				itf = x
				break
			}
		}

		// get method name and params/returns
		itfType, _ := obj.Type().Underlying().(*types.Interface)
		for i := 0; i < itfType.NumMethods(); i++ {
			x := itfType.Method(i)
			var m *Method
			for _, a := range itf.Methods {
				if a.Name == x.Name() {
					m = a
					break
				}
			}

			y := x.Type().(*types.Signature)
			m.Params = getTypeValues(y.Params())
			m.Results = getTypeValues(y.Results())
		}
	}
}

func parseImports(pkg *Package, imports []*ast.ImportSpec) {
	for _, spec := range imports {
		imp := spec.Path.Value
		// TODO must fix package name conflict
		if imp[0] == '"' {
			i := strings.LastIndex(imp, "/")
			if i == -1 {
				i = 1
			} else {
				i++
			}
			pkg.Imports[imp[i:len(imp)-1]] = imp[1 : len(imp)-1]
		} else {
			i := strings.Index(imp, " ")
			pkg.Imports[imp[:i]] = imp[i+1 : len(imp)-1]
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

func getTypeValues(tuple *types.Tuple) (vts []*VarType) {
	for i := 0; i < tuple.Len(); i++ {
		x := tuple.At(i)
		vt := &VarType{
			Var:  x.Name(),
			Deep: true,
		}
		parseType(x.Type(), vt)

		vts = append(vts, vt)
	}

	return vts
}

func parseType(t types.Type, vt *VarType) {
	tt := t.String()
	log.JSON(tt, vt)

	// TODO not deep use deep, but no reverse
	k := fmt.Sprintf("%s%t", tt, vt.Deep)
	if v, ok := parsed[k]; ok {
		tmp := *v
		tmp.Var = vt.Var
		*vt = tmp
		return
	}

	switch t := t.(type) {
	case *types.Named:
		if t.Obj() != nil {
			vt.Name = t.Obj().Name()
			if vt.Name == "error" {
				return
			}
			vt.Path = t.Obj().Pkg().Path()
			vt.Pkg = t.Obj().Pkg().Name()
		}

		if tt == "database/sql.Tx" || tt == "time.Time" {
			vt.Deep = false
			return
		}
		parseType(t.Underlying(), vt)

	case *types.Basic:
		if vt.Name != "" {
			vt.Alias = t.Name()
		} else {
			vt.Name = t.Name()
		}

	case *types.Pointer:
		parseType(t.Elem(), vt)
		vt.Pointer = "*"

	case *types.Array:
		parseType(t.Elem(), vt)
		vt.Array = fmt.Sprintf("[%d]", t.Len())

	case *types.Slice:
		parseType(t.Elem(), vt)
		vt.Slice = "[]"

	case *types.Map:
		vt.Name = "map"
		vt.Key = t.Key().String()
		vt.Value = t.Elem().String()

	case *types.Struct:
		if !vt.Deep {
			return
		}
		parseStruct(t, vt)

	default:
		log.Warnf("unimplement %#v", t)
	}

	tmp := *vt
	parsed[k] = &tmp
}

func parseStruct(t *types.Struct, x *VarType) {
	for i := 0; i < t.NumFields(); i++ {
		f := t.Field(i)
		vt := &VarType{
			Var: f.Name(),
		}
		parseType(f.Type(), vt)

		vt.Tag = getLightTag(t.Tag(i))
		setTag(vt)
		x.Fields = append(x.Fields, vt)
	}
}

func setTag(vt *VarType) {
	if vt.Var == "" || vt.Tag != "" {
		return
	}
	last := 0
	for i := 1; i < len(vt.Var); i++ {
		if vt.Var[i] >= 'A' && vt.Var[i] <= 'Z' {
			vt.Tag += vt.Var[last+1:i] + "_" + strings.ToLower(vt.Var[i:i+1])
			last = i
		}
	}
	vt.Tag = strings.ToLower(vt.Var[:1]) + vt.Tag + vt.Var[last+1:]
}

var lightRegexp = regexp.MustCompile(`light:"(.+?)"`)

func getLightTag(tag string) string {
	if tag == "" {
		return ""
	}
	m := lightRegexp.FindStringSubmatch(tag)
	if len(m) > 0 {
		return m[1]
	}
	return ""
}
