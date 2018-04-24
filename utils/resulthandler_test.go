package utils

import "testing"

func TestStripNewlineChar(t *testing.T) {
	tests := []struct {
		in  []byte
		out string
	}{
		{
			[]byte("Foo\r\n"),
			"Foo",
		},
		{
			[]byte("Foo Bar"),
			"Foo Bar",
		},
		{
			[]byte("Foo Bar\r"),
			"Foo Bar",
		},
		{
			[]byte("Foo Bar\n"),
			"Foo Bar",
		},
	}
	for _, test := range tests {
		result := StripNewlineChar(test.in)
		if result != test.out {
			t.Errorf("Epected: %s, got: %s", test.out, result)
		}
	}
}
