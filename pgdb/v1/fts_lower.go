package v1

import (
	"strings"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/filters/mapper"
)

var lowerCaseFilter = mapper.NewFilter(func(t *jargon.Token) *jargon.Token {
	if t.IsPunct() || t.IsSpace() {
		return t
	}
	v := t.String()
	next := strings.ToLower(v)
	if next == v {
		return t
	}
	return jargon.NewToken(next, true)
})
