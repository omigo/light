package null

import (
	"testing"
	"time"
)

func TestNullTimeZero(t *testing.T) {
	str := "0000-00-00 00:00:00"

	var nt NullTime

	err := nt.UnmarshalJSON([]byte(str))
	if err != nil {
		t.Error(err)
	}

	t.Log(nt.Time)

	t.Log(time.Unix(0, 0))
}
