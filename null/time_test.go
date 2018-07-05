package null

import (
	"testing"
	"time"
)

func TestNullTimeZero(t *testing.T) {
	str := `"0000-00-00 00:00:00"`

	var ckDateTime ClickHouseTime

	err := ckDateTime.UnmarshalJSON([]byte(str))
	if err != nil {
		t.Error(err)
	}

	t.Log(ckDateTime.Time)

	t.Log(time.Unix(0, 0))
}
