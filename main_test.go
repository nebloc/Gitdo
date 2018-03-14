package main

import (
	"flag"
	"testing"
)

func init() {
	HandleFlags()
}

func TestHandleFlags(t *testing.T) {
	flag.Set("v", "true")

	t.Log(verboseLogFlag)
}
