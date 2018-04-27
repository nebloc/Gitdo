package utils

import "testing"

func TestStripNewlineChar(t *testing.T) {
	tt := []struct {
		name string
		in   []byte
		out  string
	}{
		{
			"carriage return then line feed",
			[]byte("Foo\r\n"),
			"Foo",
		},
		{
			"No new line characters",
			[]byte("Foo Bar"),
			"Foo Bar",
		},
		{
			"carriage return only",
			[]byte("Foo Bar\r"),
			"Foo Bar",
		},
		{
			"line feed only",
			[]byte("Foo Bar\n"),
			"Foo Bar",
		},
		{
			"line feed then carriage return",
			[]byte("Foo Bar\n"),
			"Foo Bar",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := StripNewlineByte(tc.in)
			if result != tc.out {
				t.Fatalf("Epected: %s; got: %s", tc.out, result)
			}
		})
	}
}
