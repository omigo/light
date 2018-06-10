package goparser

import (
	"go/types"
	"regexp"
	"strings"
)

type Profile struct {
	TypeName string

	PkgName string
	PkgPath string
	Alias   string

	BasicKind types.BasicKind
	Array     bool
	Slice     bool
	Pointer   bool
	Struct    bool

	Tx bool

	Fields []*Variable `json:"-"`
}

func NewProfile(t types.Type, cache map[string]*Profile, deep bool) *Profile {
	p := &Profile{}

	str := t.String()

	switch str {
	case "*database/sql.Tx":
		p.Tx = true

	case "error":
		// log.Debug(str)

	default:
		p.parseType(t, cache, deep)
	}

	return p
}

func (p *Profile) parseType(t types.Type, cache map[string]*Profile, deep bool) {
	switch v := t.(type) {
	case *types.Basic:
		p.TypeName = v.Name()
		p.BasicKind = v.Kind()

	case *types.Map:
		panic("unsupported type " + v.String())

	case *types.Named:
		if obj := v.Obj(); obj != nil {
			p.TypeName = obj.Name()
			if pkg := obj.Pkg(); pkg != nil {
				p.PkgName = pkg.Name()
				p.PkgPath = pkg.Path()
			}
			if s, ok := v.Underlying().(*types.Struct); ok {
				p.Struct = true
				if deep {
					p.parseStruct(s, cache)
				}
			} else {
				p.parseType(v.Underlying(), cache, deep)
				tstr := v.Obj().Type().String()
				if p.PkgPath != "" && strings.HasPrefix(tstr, p.PkgPath) {
					p.Alias = tstr[len(p.PkgPath):]
					p.Alias = strings.TrimPrefix(p.Alias, ".")
				}
			}
		}

	case *types.Pointer:
		p.Pointer = true
		p.parseType(v.Elem(), cache, deep)

	case *types.Array:
		p.Array = true
		p.parseType(v.Elem(), cache, deep)

	case *types.Slice:
		p.Slice = true
		p.parseType(v.Elem(), cache, deep)

	case *types.Struct:
		p.Struct = true
		if deep {
			p.parseStruct(v, cache)
		}

	case *types.Chan, *types.Interface, *types.Signature, *types.Tuple:
		panic("unsupported type " + v.String())
	}
}

func (p *Profile) parseStruct(s *types.Struct, cache map[string]*Profile) {
	for i := 0; i < s.NumFields(); i++ {
		alias, cmds := parseTags(s.Tag(i))
		v := NewVariableTag(s.Field(i), alias, cmds)
		p.Fields = append(p.Fields, v)
	}
}

var tagRegexp = regexp.MustCompile(`(.+):"(.+)"`)

func parseTags(tag string) (alias string, cmds []string) {
	// Username string `json:"username" light:"uname,nullable"`

	groups := tagRegexp.FindAllStringSubmatch(tag, -1)
	for _, m := range groups {
		if m[1] != "light" {
			continue
		}
		vs := strings.Split(m[2], ",")
		if len(vs) == 0 {
			return "", nil
		} else if len(vs) == 1 {
			return vs[0], nil
		} else {
			return vs[0], vs[1:]
		}
	}

	return "", nil
}

func (p *Profile) FullTypeName() string {
	var name string
	if p.Slice {
		name += "[]"
	}
	if p.Pointer {
		name += "*"
	}
	if p.PkgName != "" {
		name += p.PkgName + "."
	}
	return name + p.TypeName
}

func (p *Profile) FullElemTypeName() string {
	var name string
	if p.PkgName != "" {
		name += p.PkgName + "."
	}
	return name + p.TypeName
}
