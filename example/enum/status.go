package enum

type Status int32

const (
	StatusUnknow Status = iota
	StatusNormal
	StatusDeleted
)
