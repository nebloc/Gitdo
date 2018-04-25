package main

import (
	"github.com/nebloc/gitdo/cmd"
)

var version string

func main() {
	cmd.New(version).Execute()
}
