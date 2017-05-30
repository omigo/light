package main

import (
	"go/ast"
	"go/parser"
	"go/types"

	"github.com/arstd/log"
	"golang.org/x/tools/go/loader"
)

var parsed = map[string]*VarType{}

func parseGoFileByLoader(pkg *Package) {
	defer func() { parsed = nil }()

	conf := loader.Config{
		ParserMode:          parser.ParseComments,
		TypeCheckFuncBodies: func(path string) bool { return false },
	}
	conf.CreateFromFilenames("arstd/light", pkg.Source)
	prog, err := conf.Load()
	if err != nil {
		log.Panic(err)
	}

	pkgInfos := prog.InitialPackages()
	info := pkgInfos[0]

	pkg.Name = info.Pkg.Name()
	parseImports(pkg, info.Files[0].Imports)

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

		itf := &Interface{
			Name: typeSpec.Name.Name,
		}
		pkg.Interfaces = append(pkg.Interfaces, itf)

		// get method name and doc
		for _, x := range interfaceType.Methods.List {
			m := &Method{
				Name: x.Names[0].Name,
				Doc:  getDoc(x.Doc),
			}
			itf.Methods = append(itf.Methods, m)
		}

		// get method name and params/returns
		itfType, _ := v.Type().Underlying().(*types.Interface)
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
