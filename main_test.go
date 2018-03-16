package main

import (
	"flag"
	"testing"
)

func init() {
	SetArgs()
}

func TestSetArgs(t *testing.T) {
	flag.Set("v", "true")

	t.Log(verboseLogFlag)
}
