package v1

import (
	"strings"
	"unicode/utf8"
)

// ToValidUTF8 is a fork of the Go standard library's strings.ToValidUTF8, differing in the following ways:
// * It uses utf8.RuneError as the replacement rune.
// * It replaces each invalid rune instead of replacing all adjacent invalid runes with a single replacement.
// * It drops 0-bytes from the resulting string.
func ToValidUTF8(s string) string {
	var b strings.Builder

	// fastpath: scan until we find the first null or invalid rune, then allocate and copy.
	for i, c := range s {
		if c == 0 {
			b.Grow(len(s) - 1)
			b.WriteString(s[:i])
			s = s[i:]
			break
		}

		_, wid := utf8.DecodeRuneInString(s[i:])
		if wid == 1 {
			b.Grow(len(s) + 3) // utf8.RuneLen(utf8.RuneError)
			b.WriteString(s[:i])
			s = s[i:]
			break
		}
	}
	if b.Cap() == 0 { // didn't call b.Grow above
		return s
	}

	// else, we need to process the input
	for i := 0; i < len(s); {
		c := s[i]

		// skip nulls
		if c == 0 {
			i++
			continue
		}

		// characters below utf8.RuneSelf are valid single bytes and don't need replaced
		if c < utf8.RuneSelf {
			i++
			b.WriteByte(c)
			continue
		}

		// at this point, any width of 1 is an error.
		_, wid := utf8.DecodeRuneInString(s[i:])
		if wid == 1 {
			i++
			b.WriteRune(utf8.RuneError)
			continue
		}
		b.WriteString(s[i : i+wid])
		i += wid
	}

	return b.String()
}

func SanitizeString(s string) string {
	return ToValidUTF8(s)
}
