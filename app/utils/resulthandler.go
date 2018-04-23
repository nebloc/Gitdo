package utils

import "strings"

// stripNewLineChar takes a byte array (usually from an exec.Command run) and strips the newline characters, returning
// a string
func StripNewlineByte(bytes []byte) string {
	return StripNewlineString(string(bytes))
}

func StripNewlineString(str string) string {
	// Strip line feed
	if strings.HasSuffix(str, "\n") {
		str = str[:len(str)-1]
	}
	// Strip carriage return
	if strings.HasSuffix(str, "\r") {
		str = str[:len(str)-1]
	}

	return str
}
