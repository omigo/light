package null

import (
	"testing"
	"unsafe"
)

func TestIntSize(t *testing.T) {
	var a int
	t.Log(unsafe.Sizeof(a))
}
