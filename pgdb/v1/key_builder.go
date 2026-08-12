package v1

import "strings"

// KeyBuilder builds pk/sk key strings for generated Record() and PKSK()
// code. It wraps strings.Builder and silently drops NUL (U+0000) bytes as
// parts are written: Postgres text values can never contain NUL, and under
// pgx's simple protocol an interpolated NUL truncates the query wire
// message mid-statement, so the server rejects the whole statement with
// 08P01 "invalid message format". Dropping NUL here keeps keys consistent
// with plain string columns, which strip NUL at Record() time.
type KeyBuilder struct {
	sb strings.Builder
}

// Reset resets the builder to be empty.
func (b *KeyBuilder) Reset() { b.sb.Reset() }

// String returns the accumulated string.
func (b *KeyBuilder) String() string { return b.sb.String() }

// WriteString appends s, minus any NUL bytes, to the builder. It returns
// the number of bytes actually written and a nil error, matching the
// strings.Builder signature.
func (b *KeyBuilder) WriteString(s string) (int, error) {
	n := 0
	for {
		i := strings.IndexByte(s, 0)
		if i < 0 {
			w, _ := b.sb.WriteString(s)
			return n + w, nil
		}
		w, _ := b.sb.WriteString(s[:i])
		n += w
		s = s[i+1:]
	}
}
