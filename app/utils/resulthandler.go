package utils

import "strings"

// stripNewLineChar takes a byte array (usually from an exec.Command run) and strips the newline characters, returning
// a string
func StripNewlineChar(orig []byte) string {
	var newStr string
	// Strip line feed
	if strings.HasSuffix(string(orig), "\n") {
		newStr = string(orig)[:len(orig)-1]
	}
	// Strip carriage return
	if strings.HasSuffix(newStr, "\r") {
		newStr = newStr[:len(newStr)-1]
	}
	return newStr
}
