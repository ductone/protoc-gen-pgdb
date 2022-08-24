package v1

import (
	"bytes"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/uax29/graphemes"
	"github.com/ductone/protoc-gen-pgdb/internal/slice"
)

func edgegramStream(n int) jargon.Filter {
	return func(incoming *jargon.TokenStream) *jargon.TokenStream {
		rg := &edgegram{size: n, incoming: incoming}
		return jargon.NewTokenStream(rg.next)
	}
}

type edgegram struct {
	size     int
	incoming *jargon.TokenStream
	pending  []*jargon.Token
}

func (t *edgegram) shiftPending() (*jargon.Token, error) {
	if len(t.pending) >= 1 {
		rv := t.pending[0]
		t.pending = t.pending[1:]
		return rv, nil
	}
	return nil, nil
}

func (t *edgegram) next() (*jargon.Token, error) {
	if tk, err := t.shiftPending(); tk != nil || err != nil {
		return tk, err
	}

	token, err := t.incoming.Next()
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, nil
	}

	if token.IsPunct() || token.IsSpace() {
		return token, nil
	}
	v := token.String()
	t.pending = make([]*jargon.Token, 0, (len(v)+3)/3)

	segments := graphemes.NewSegmenter([]byte(v))
	sg := graphemes.SegmentAll([]byte(v))
	offset := 0
	grams := make([]string, 0, (len(sg)+3)/3)
	for i := 1; i <= len(sg); i++ {
		k := string(bytes.Join(sg[0:i], []byte{}))
		grams = append(grams, k)
		if i > t.size {
			offset = i - t.size
			k := string(bytes.Join(sg[offset:i], []byte{}))
			grams = append(grams, k)
		}
	}
	grams = slice.Unique(grams)
	slice.Sort(grams)
	for _, k := range grams {
		t.pending = append(t.pending, jargon.NewToken(k, true))
	}

	if err := segments.Err(); err != nil {
		return nil, err
	}

	return t.shiftPending()
}
