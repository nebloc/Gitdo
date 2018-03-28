package main

import (
	"fmt"
	"github.com/fatih/color"
)

func Warn(out string) {
	color.Yellow(out)
}
func Warnf(format string, out ...interface{}) {
	color.Yellow(fmt.Sprintf(format, out...))
}

func Danger(out string) {
	color.Red(out)
}
func Dangerf(format string, out ...interface{}) {
	color.Red(fmt.Sprintf(format, out...))
}

func Highlight(out string) {
	color.Cyan(out)
}

func Highlightf(format string, out ...interface{}) {
	color.Cyan(fmt.Sprintf(format, out...))
}
