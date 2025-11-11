// ascii.go
// ---------
// Converts Unicode filenames to ASCII-safe equivalents using NFKD normalization.

package main

import (
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func cleanASCII(s string) string {
	if s == "" {
		return s
	}
	decomp := norm.NFKD.String(s)
	out := make([]rune, 0, len(decomp))
	for _, r := range decomp {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		if r < 128 {
			if r >= 32 {
				out = append(out, r)
			}
			continue
		}
		switch r {
		case 'ß':
			out = append(out, 's', 's')
		case 'Æ', 'Ǽ', 'æ', 'ǽ':
			out = append(out, 'a', 'e')
		case 'Œ', 'œ':
			out = append(out, 'o', 'e')
		case 'Đ', 'đ', 'Ð', 'ð':
			out = append(out, 'd')
		case 'Ł', 'ł':
			out = append(out, 'l')
		case '₫':
			out = append(out, 'd')
		case '–', '—':
			out = append(out, '-')
		case '“', '”', '„', '«', '»', '′', '″':
			out = append(out, '\'')
		case '·', '•', '∙':
			out = append(out, '-')
		default:
			out = append(out, '_')
		}
	}
	return string(out)
}
