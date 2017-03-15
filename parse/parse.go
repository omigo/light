package parse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/types"
	"strings"

	"github.com/arstd/light/domain"
	"github.com/arstd/light/util"
	"github.com/arstd/log"

	"golang.org/x/tools/go/loader"
)

func ParseGoFile(file string) (pkg *domain.Package) {
	pkg = &domain.Package{Source: file}

	conf := loader.Config{
		ParserMode:          parser.ParseComments,
		TypeCheckFuncBodies: func(path string) bool { return false },
	}
	conf.CreateFromFilenames("arstd/light", file)
	prog, err := conf.Load()
	util.CheckError(err)

	pkgInfos := prog.InitialPackages()
	info := pkgInfos[0]

	pkg.Path = info.Pkg.Path()
	pkg.Name = info.Pkg.Name()
	pkg.Imports = parseImports(info.Files[0].Imports)

	for k, v := range info.Defs {
		if k.Obj == nil || k.Obj.Kind != ast.Typ {
			continue
		}
		typeSpec, ok := k.Obj.Decl.(*ast.TypeSpec)
		if !ok {
			continue
		}
		interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
		if !ok {
			continue
		}

		itf := &domain.Interface{
			Name: typeSpec.Name.Name,
		}
		pkg.Interfaces = append(pkg.Interfaces, itf)

		itfType, _ := v.Type().Underlying().(*types.Interface)
		for i, x := range interfaceType.Methods.List {
			m := &domain.Method{
				Name: x.Names[0].Name,
				Doc:  getDoc(x.Doc),
			}
			itf.Methods = append(itf.Methods, m)

			y := itfType.Method(i).Type().(*types.Signature)
			m.Params = getTypeValues(y.Params())
			m.Results = getTypeValues(y.Results())
			checkResultsVar(m)
		}
	}

	return pkg
}

func parseImports(imports []*ast.ImportSpec) (ret map[string]string) {
	ret = map[string]string{}
	for _, spec := range imports {
		imp := spec.Path.Value
		// TODO must fix package name conflict
		if imp[0] == '"' {
			i := strings.LastIndex(imp, "/")
			ret[imp[i+1:len(imp)-1]] = imp[1 : len(imp)-1]
		} else {
			i := strings.Index(imp, " ")
			ret[imp[:i]] = imp[i+1 : len(imp)-1]
		}
	}
	return ret
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

func getTypeValues(tuple *types.Tuple) (vts []*domain.VarType) {
	for i := 0; i < tuple.Len(); i++ {
		x := tuple.At(i)
		vt := &domain.VarType{
			Var:  x.Name(),
			Deep: true,
		}
		parseType(x.Type(), vt)

		vts = append(vts, vt)
	}

	return vts
}

var parsed = map[string]*domain.VarType{}

func parseType(t types.Type, vt *domain.VarType) {

	tt := t.String()
	k := tt + fmt.Sprint(vt.Deep)
	// log.Debug(k)
	if v, ok := parsed[k]; ok {
		tmp := *v
		tmp.Var = vt.Var
		tmp.Deep = vt.Deep
		if !tmp.Deep {
			tmp.Fields = nil
		}
		*vt = *(&tmp)
		return
	}
	parsed[k] = vt

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

		// log.Warnf("named %#v", t.String())

		if tt == "database/sql.Tx" || tt == "time.Time" {
			return
		}

		if !vt.Deep {
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
		vt.Pointer = "*"
		parseType(t.Elem(), vt)

	case *types.Array:
		vt.Array = fmt.Sprintf("[%d]", t.Len())
		parseType(t.Elem(), vt)

	case *types.Slice:
		vt.Slice = "[]"
		parseType(t.Elem(), vt)

	case *types.Map:
		vt.Name = "map"
		vt.Key = t.Key().String()
		vt.Elem = t.Elem().String()

	case *types.Struct:
		if !vt.Deep {
			return
		}
		parseStruct(t, vt)

	default:
		log.Warnf("unimplement %#v", t)
	}
}

func parseStruct(t *types.Struct, x *domain.VarType) {
	for i := 0; i < t.NumFields(); i++ {
		f := t.Field(i)
		vt := &domain.VarType{
			Var: f.Name(),
		}

		// log.Infof("%#v", f.String())
		parseType(f.Type(), vt)
		x.Fields = append(x.Fields, vt)
	}
}

func checkResultsVar(m *domain.Method) {
	for _, vt := range m.Results {
		if vt.Var == "" {
			if vt.Name == "error" {
				vt.Var = "err"
			} else {
				vt.Var = "x"
				if vt.Pkg != "" {
					vt.Var += vt.Pkg[:1]
				}
				if vt.Name != "" {
					vt.Var += strings.ToLower(vt.Name[:1])
				}
				if vt.Slice != "" {
					vt.Var += "s"
				}
			}
			for _, v := range m.Params {
				if v.Var == "vt.Var" {
					vt.Var = "x" + vt.Var
					break
				}
			}
		}
	}
}
