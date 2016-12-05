package enum

type State string

const (
	StateUnknow  State = "unknow"
	StateNormal        = "normal"
	StateDeleted       = "deleted"
)

func (s State) String() string {
	return string(s)
}
