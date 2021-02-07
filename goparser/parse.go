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
	"strconv"
	"strings"

	"github.com/omigo/log"
)

type Interface struct {
	Source  string
	Log     bool
	Timeout int64

	Package string            // itf
	Imports map[string]string // database/sql => sql
	Name    string            // IUser

	VarName   string
	StoreName string

	Methods []*Method

	// full-type-name : type-profile
	Cache map[string]*Profile
}

func Parse(filename string, src interface{}) (*Interface, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		log.Panic(err)
	}
	// ast.Print(fset, f)

	itf := &Interface{
		Source:  filename,
		Package: f.Name.Name,
		Imports: map[string]string{},
	}

	goBuild(filename)

	extractDocs(itf, f, fset)

	extractTypes(itf, f, fset)

	// log.JsonIndent(itf)

	itf.makeCache()

	return itf, nil
}

func goBuild(src string) {
	cmd := exec.Command("go", "build", src)
	out, err := cmd.CombinedOutput()
	if bytes.HasSuffix(out, []byte("command-line-arguments\n")) {
		fmt.Printf("%s", out[:len(out)-23])
	} else {
		fmt.Printf("%s", out)
	}
	if err != nil {
		panic(err)
	}
}

func extractDocs(itf *Interface, f *ast.File, fset *token.FileSet) {
	for _, decl := range f.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			switch genDecl.Tok {
			case token.IMPORT:
				for _, spec := range genDecl.Specs {
					if importSpec, ok := spec.(*ast.ImportSpec); ok {
						path, err := strconv.Unquote(importSpec.Path.Value)
						if err != nil {
							panic(importSpec.Path.Value + " " + err.Error())
						}
						if importSpec.Name != nil {
							itf.Imports[path] = importSpec.Name.Name
						} else {
							itf.Imports[path] = ""
						}
					}
				}

			case token.TYPE:
				for _, spec := range genDecl.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
							if itf.Name != "" {
								panic("one file must contains one interface only")
							}

							itf.Name = typeSpec.Name.Name
							for _, method := range interfaceType.Methods.List {
								m := NewMethod(itf, method.Names[0].Name, getDoc(method.Doc))
								itf.Methods = append(itf.Methods, m)
							}
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

func extractTypes(itf *Interface, f *ast.File, fset *token.FileSet) {
	info := types.Info{Defs: make(map[*ast.Ident]types.Object)}
	conf := types.Config{Importer: importer.For("source", nil)}
	_, err := conf.Check(itf.Package, fset, []*ast.File{f}, &info)
	log.Fataln(err)

	for k, obj := range info.Defs {
		if k.Obj != nil {
			if k.Name == itf.Name {
				if k.Obj.Kind == ast.Typ {
					// get method name and params/returns
					if itfType, ok := obj.Type().Underlying().(*types.Interface); ok {
						for i := 0; i < itfType.NumMethods(); i++ {
							x := itfType.Method(i)
							m := getMethodByName(itf, x.Name())
							y := x.Type().(*types.Signature)
							m.Params = NewParams(y.Params())
							m.Results = NewResults(y.Results())
						}
					}
				}
			} else {
				if tn, ok := obj.Type().(*types.Named); ok {
					if itf.Name == tn.Obj().Name() {
						itf.VarName = k.Name
					}
				}
			}
		}
	}
}

func getMethodByName(s *Interface, name string) *Method {
	for _, a := range s.Methods {
		if a.Name == name {
			return a
		}
	}
	return nil
}

func (itf *Interface) makeCache() {
	itf.Cache = map[string]*Profile{}

	for _, method := range itf.Methods {
		for _, param := range method.Params.List {
			key := param.Type.String()
			profile, ok := itf.Cache[key]
			if !ok {
				profile = NewProfile(param.Type, itf.Cache, true)
				itf.Cache[key] = profile
			}

			for _, f := range profile.Fields {
				if f.PkgPath != "" {
					itf.Imports[f.PkgPath] = ""
				}

				// field 是一个变量，在不同的方法中，名字不一样，所以不能公用
				field := new(Variable)
				*field = *f

				k := field.Type.String()
				p, ok := itf.Cache[k]
				if !ok {
					p = NewProfile(field.Type, itf.Cache, false)
					itf.Cache[k] = p
				}
				field.Profile = p
				field.Parent = param

				method.Params.Names[field.Name] = field
				method.Params.Names[underLower(field.Name)] = field
				method.Params.Names[param.Name+"."+field.Name] = field
				if field.TagAlias != "" {
					method.Params.Names[field.TagAlias] = field
				}
				if profile.Slice {
					if param.Name[len(param.Name)-1] == 's' {
						elem := param.Name[:len(param.Name)-1]
						method.Params.Names[elem+"."+field.Name] = field
					}
				}
			}
			*param.Profile = *profile
			method.Params.Names[param.Name] = param
			method.Params.Names[underLower(param.Name)] = param
		}
		for _, result := range method.Results.List {
			result.Name = ""
			key := result.Type.String()
			profile, ok := itf.Cache[key]
			if !ok {
				profile = NewProfile(result.Type, itf.Cache, true)
				itf.Cache[key] = profile
			}

			for _, f := range profile.Fields {
				if f.PkgPath != "" {
					itf.Imports[f.PkgPath] = ""
				}

				field := new(Variable)
				*field = *f

				k := field.Type.String()
				p, ok := itf.Cache[k]
				if !ok {
					p = NewProfile(field.Type, itf.Cache, false)
					itf.Cache[k] = p
				}
				field.Profile = p
				field.Parent = result
				if field.Name == "" {
					log.JsonIndent(profile)
					panic("unreachable code")
				}
				method.Results.Names[field.Name] = field
				method.Results.Names[underLower(field.Name)] = field
				if result.Name != "" {
					method.Results.Names[result.Name+"."+field.Name] = field
				}
				if field.TagAlias != "" {
					method.Results.Names[field.TagAlias] = field
				}
			}
			result.Profile = profile
			if result.Name != "" {
				method.Results.Names[result.Name] = result
			}
		}
	}
}

func underLower(field string) string {
	var buf bytes.Buffer
	for i, v := range field {
		if v >= 'A' && v <= 'Z' {
			if i != 0 {
				buf.WriteByte('_')
			}
			buf.WriteRune(v + 32)
		} else {
			buf.WriteRune(v)
		}
	}
	return buf.String()
}

func upperCamelCase(field string) string {
	var buf bytes.Buffer
	var upper bool = true
	for _, v := range field {
		if v == '_' {
			upper = true
		} else if upper {
			buf.WriteRune(v - 32)
			upper = false
		} else {
			buf.WriteRune(v)
		}
	}
	return buf.String()
}
