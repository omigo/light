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
	// ast.Print(fset, f)

	pkg.Name = f.Name.Name
	parseImports(pkg, f.Imports)

	// log.JSONIndent(pkg)

	parseComments(pkg, f)

	// printer.Fprint(os.Stdout, fset, f)

	parseTypes(pkg, fset, f)
}

func goBuild(goFile string) {
	// log.Debugf("go build -i -v  %s", goFile)
	cmd := exec.Command("go", "build", "-i", "-v", goFile)
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
			m.Params = getTypeValues(pkg, y.Params())
			m.Results = getTypeValues(pkg, y.Results())
		}
	}
}

func parseImports(pkg *Package, imports []*ast.ImportSpec) {
	for _, spec := range imports {
		var short string
		if spec.Name != nil {
			short = spec.Name.Name
		}
		path := spec.Path.Value
		path = path[1 : len(path)-1]

		pkg.Imports[path] = short
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

func getTypeValues(p *Package, tuple *types.Tuple) (vts []*VarType) {
	for i := 0; i < tuple.Len(); i++ {
		x := tuple.At(i)
		vt := &VarType{
			Var:  x.Name(),
			Deep: true,
		}
		parseType(p, x.Type(), vt)

		vts = append(vts, vt)
	}

	return vts
}

func parseType(p *Package, t types.Type, vt *VarType) {
	tt := t.String()
	// log.JSON(tt, vt)

	// TODO not deep use deep, but no reverse
	k := fmt.Sprintf("%s%t", tt, vt.Deep)
	if v, ok := parsed[k]; ok {
		if vt.Name != "" {
			if v.Alias != "" {
				vt.Alias = v.Alias
			} else {
				vt.Alias = v.Name
			}
			return
		}
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
			vt.Pkg = p.Imports[vt.Path]
			if vt.Pkg == "" {
				vt.Pkg = t.Obj().Pkg().Name()
			}
		}

		if tt == "database/sql.Tx" || tt == "time.Time" {
			vt.Deep = false
			return
		}

		parseType(p, t.Underlying(), vt)

	case *types.Basic:
		log.Debug(vt, t.Name())
		if vt.Name != "" {
			vt.Alias = t.Name()
		} else {
			vt.Name = t.Name()
		}

	case *types.Pointer:
		parseType(p, t.Elem(), vt)
		vt.Pointer = "*"

	case *types.Array:
		parseType(p, t.Elem(), vt)
		vt.Array = fmt.Sprintf("[%d]", t.Len())

	case *types.Slice:
		parseType(p, t.Elem(), vt)
		vt.Slice = "[]"

	case *types.Map:
		vt.Name = "map"
		vt.Key = t.Key().String()
		vt.Value = t.Elem().String()

	case *types.Struct:
		if !vt.Deep {
			return
		}
		parseStruct(p, t, vt)

	default:
		log.Warnf("unimplement %#v", t)
	}

	tmp := *vt
	var hasPath string
	if tmp.Path != "" {
		hasPath = "."
	}
	k = tmp.Slice + tmp.Pointer + tmp.Path + hasPath + tmp.Name + fmt.Sprintf("%t", tmp.Deep)
	parsed[k] = &tmp
}

func parseStruct(p *Package, t *types.Struct, x *VarType) {
	for i := 0; i < t.NumFields(); i++ {
		f := t.Field(i)
		vt := &VarType{
			Var: f.Name(),
		}
		parseType(p, f.Type(), vt)

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
