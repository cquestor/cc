package cc

import (
	"unicode"
)

func camel2Case(s string) string {
	var output []rune
	for i, r := range s {
		if i == 0 {
			output = append(output, unicode.ToLower(r))
		} else {
			if unicode.IsUpper(r) {
				output = append(output, '_')
			}
			output = append(output, unicode.ToLower(r))
		}
	}
	return string(output)
}
