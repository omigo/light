package domain

type MethodKind string

const (
	Insert MethodKind = "insert"
	Update            = "update"
	Delete            = "delete"
	Get               = "get"
	Count             = "count"
	List              = "list"
	Page              = "page"
)

type Package struct {
	Path string
	Name string

	Imports map[string]string

	Interfaces []*Interface
}

type Interface struct {
	Name string

	Methods []*Func
}

type Func struct {
	Name string
	Doc  string

	Params  []*VarType
	Results []*VarType
}

type VarType struct {
	// ms []*domain.Model
	Var  string `json:"Var,omitempty"`  //  ms
	Type string `json:"Type,omitempty"` //  []*domain.Model

	Path    string `json:"Path,omitempty"`    //  github.com/arstd/light/example/domain
	Array   string `json:"Array,omitempty"`   //  []
	Slice   string `json:"Slice,omitempty"`   //  []
	Pointer string `json:"Pointer,omitempty"` //  *
	Pkg     string `json:"Pkg,omitempty"`     //  domain
	Name    string `json:"Name,omitempty"`    //  Model
	Alias   string `json:"Alias,omitempty"`   //  e.g. domain.State => string

	Key  string `json:"Key,omitempty"`
	Elem string `json:"Elem,omitempty"`

	Deep   bool       `json:"Deep,omitempty"` //  深入解析这个类型
	Fields []*VarType `json:"Fields,omitempty"`
}
