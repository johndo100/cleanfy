// case.go
// --------
// Provides case transformations: lower, upper, and title.

package main

import (
	"unicode"
	"unicode/utf8"
)

func toTitle(s string) string {
	var out []rune
	capNext := true
	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		s = s[size:]
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if capNext {
				out = append(out, unicode.ToTitle(r))
				capNext = false
			} else {
				out = append(out, unicode.ToLower(r))
			}
			continue
		}
		capNext = true
		out = append(out, r)
	}
	return string(out)
}
