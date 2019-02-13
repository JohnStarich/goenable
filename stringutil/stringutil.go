package stringutil

import (
	"strings"
)

// Dedent removes leading whitespace from every line, evenly
func Dedent(str string) string {
	smallestLeadingSpace := len(str)
	lines := strings.Split(str, "\n")
	for _, line := range lines {
		i := 0
		containedNonSpace := false
		for _, c := range line {
			switch c {
			case rune('\t'), rune(' '):
				i++
			default:
				containedNonSpace = true
				break
			}
		}
		if containedNonSpace && i < smallestLeadingSpace {
			smallestLeadingSpace = i
		}
	}
	if smallestLeadingSpace == 0 {
		return str
	}

	result := ""
	for _, line := range lines {
		if len(line) < smallestLeadingSpace {
			result += line
		} else {
			result += line[smallestLeadingSpace:]
		}
		result += "\n"
	}
	return result
}

// SingleQuote makes the given string into a single quoted entity for a shell
// If the string is simple, then it is returned as-is.
// If the string contains single quotes ('), then it its quotes are escaped.
func SingleQuote(str string) string {
	if !strings.Contains(str, `'`) {
		return str
	}
	return `'` + strings.Replace(str, `'`, `'\''`, -1) + `'`
}
