package agent

import (
	"reflect"
	"testing"
)

func TestConvertStrToByte(t *testing.T) {
	cases := [][]string{
		{"戶戸户"},
		{"test", "( ͡° ͜ʖ ͡°)"},
		{"test", "test", "test 👽 test"},
		{""},
	}

	for _, c := range cases {
		b, err := convertStrToByte(c)

		if err != nil {
			t.Log(err)
		}

		out := convertByteToStr(b)

		ok := reflect.DeepEqual(c, out)

		if !ok {
			t.Error("string arrays not equal")
		}

		t.Logf("%#v equal %#v", c, out)
	}
}
