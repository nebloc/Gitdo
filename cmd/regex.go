package cmd

import (
	"regexp"
)

var (
	// TODO: Create a library of regex's for use with other languages. <OaTSrQjZ>
	// todoReg is a compiled regex to match the TODO comments
	todoReg = regexp.MustCompile(
		`^[[:space:]]*(?://|#)[[:space:]]*TODO:[[:space:]]*(.*)`)
	taggedReg = regexp.MustCompile(
		`^[[:space:]]*(?://|#)[[:space:]]*TODO(?::|)[[:space:]]*(?:.*)<(.*)>`)
	looseTODOReg = regexp.MustCompile(
		`^[[:space:]]*(?://|#)[[:space:]]*TODO(?::|)[[:space:]]*(.*)`)
)

// Check regex, takes a regex, attempts to match it against a given string, and returns if it matched, and the first capture group.
func CheckRegex(reg *regexp.Regexp, str string) (string, bool) {
	match := reg.FindStringSubmatch(str)
	if len(match) < 2 {
		return "", false
	}
	return match[1], true
}
