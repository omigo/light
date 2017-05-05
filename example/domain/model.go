package domain

import (
	"time"

	"github.com/arstd/light/example/enum"
)

// Model 模型示例
type Model struct {
	Id    int
	Name  string
	Flag  bool
	Score float32

	Map   map[string]interface{}
	Time  time.Time
	Array []int64
	Slice []string `light:"slice text[]"`

	Status enum.Status
	State  enum.State

	Pointer     *Model
	StructSlice []*Model
	Uint32      uint32
}
