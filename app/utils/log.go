package utils

import (
	"fmt"
	"github.com/fatih/color"
)

// Prints to the console in a yellow colour
func Warn(out string) {
	color.Yellow(out)
}

// Formats the string and prints to the console in a yellow colour
func Warnf(format string, out ...interface{}) {
	color.Yellow(fmt.Sprintf(format, out...))
}

// Prints to the console in a red colour
func Danger(out string) {
	color.Red(out)
}

// Formats the string and prints to the console in a red colour
func Dangerf(format string, out ...interface{}) {
	color.Red(fmt.Sprintf(format, out...))
}

// Prints to the console in a blue colour
func Highlight(out string) {
	color.Cyan(out)
}

// Formats the string and prints to the console in a blue colour
func Highlightf(format string, out ...interface{}) {
	color.Cyan(fmt.Sprintf(format, out...))
}
