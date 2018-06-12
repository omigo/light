package goparser

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/arstd/log"
)

func ResultWrap(v *Variable) string         { return v.Wrap() }
func ResultTypeName(v *Variable) string     { return v.FullTypeName() }
func ResultElemTypeName(v *Variable) string { return v.FullElemTypeName() }
func LookupScanOfResults(m *Method, name string) string {
	v := m.Results.Lookup(name)
	if v == nil {
		// log.Error(fmt.Sprintf("method `%s` result varialbe `%s` not found", m.Name, name))
		return m.Results.Result.Scan()
	}

	return v.Scan()
}
func LookupValueOfParams(m *Method, name string) string {
	v := m.Params.Lookup(name)
	if v == nil {
		panic(fmt.Sprintf("method `%s` result varialbe `%s` not found", m.Name, name))
	}
	return v.Value()
}

type Variable struct {
	Name string
	*Profile

	Parent *Variable

	TagAlias string
	TagCmds  []string

	Var  *types.Var
	Type types.Type
}

func NewVariable(v *types.Var) *Variable {
	return NewVariableTag(v, "", nil)
}

func NewVariableTag(v *types.Var, tagAlias string, tagCmds []string) *Variable {
	variable := &Variable{
		Name:    v.Name(),
		Profile: new(Profile),

		TagAlias: tagAlias,
		TagCmds:  tagCmds,

		Var:  v,
		Type: v.Type(),
	}
	return variable
}

func (v *Variable) Nullable() bool {
	for _, cmd := range v.TagCmds {
		if cmd == "nullable" {
			return true
		}
	}
	return false
}

func (v *Variable) NotDefault() string {
	name := v.FullName()

	switch {
	case v.PkgPath == "github.com/arstd/light/null":
		return "!" + name + ".IsEmpty()"

	case v.PkgPath == "time" && v.TypeName == "Time":
		return "!" + name + ".IsZero()"

	case v.Pointer:
		return name + " != nil"

	case v.Struct:
		return name + " != nil"

	case v.Array:
		return name + " != nil"

	case v.Slice:
		return "len(" + name + ") != 0"

	case v.BasicKind == types.String:
		return name + ` != ""`

	case v.BasicKind == types.Bool:
		return "!" + name

	case v.BasicKind >= types.Int && v.BasicKind <= types.Uint64:
		return name + ` != 0`

	default:
		log.JsonIndent(v)
		panic("unimplement not default for variable " + v.PkgPath + "." + v.TypeName)
	}
}

func (v *Variable) FullName() string {
	var name string

	if v.Parent != nil {
		if v.Parent.Name == "" {
			name += "xu."
		} else {
			name += v.Parent.Name + "."
		}
	}
	if v.Name == "" {
		return name + "xu"
	}

	return name + v.Name
}

func (v *Variable) Scan() string {
	name := v.FullName()
	switch {
	case v.PkgPath == "github.com/arstd/light/null":
		return "&" + name
	case v.Pointer:
		return "&" + name
	case v.Nullable():
		return fmt.Sprintf("null.%s%s(&%s)", strings.ToUpper(v.TypeName[:1]), v.TypeName[1:], name)
	default:
		return "&" + name
	}
}

func (v *Variable) Wrap() string {
	name := v.FullName()
	if v.PkgPath == "github.com/arstd/light/null" {
		return name
	}
	name = fmt.Sprintf("null.%s%s(&%s)", strings.ToUpper(v.TypeName[:1]), v.TypeName[1:], name)
	return name
}

func (v *Variable) Value() string {
	name := v.FullName()
	switch {
	case v.PkgPath == "github.com/arstd/light/null":
		return name
	case v.Pointer:
		return name
	case v.Nullable():
		return fmt.Sprintf("null.%s%s(&%s)", strings.ToUpper(v.TypeName[:1]), v.TypeName[1:], name)
	default:
		return name
	}
}

func (v *Variable) Define() string {
	return v.Name + " " + v.FullTypeName()
}
