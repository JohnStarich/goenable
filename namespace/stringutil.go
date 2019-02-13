package namespace

import (
	"strings"
)

func dedent(str string) string {
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

func singleQuote(str string) string {
	if !strings.Contains(str, `'`) {
		return str
	}
	return `'` + strings.Replace(str, `'`, `'\''`, -1) + `'`
}
