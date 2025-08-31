package v1

import (
	"bytes"
	"sort"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/uax29/graphemes"
)

func edgegramStream(n int) jargon.Filter {
	return func(incoming *jargon.TokenStream) *jargon.TokenStream {
		rg := &edgegram{size: n, incoming: incoming}
		return jargon.NewTokenStream(rg.next)
	}
}

type byteRange struct {
	start int
	end   int
}

type edgegram struct {
	size     int
	incoming *jargon.TokenStream

	vb     []byte
	grams  []byteRange
	cursor int
}

func (t *edgegram) nextFromGrams() (*jargon.Token, error) {
	if t.cursor < len(t.grams) {
		r := t.grams[t.cursor]
		t.cursor++
		return jargon.NewToken(string(t.vb[r.start:r.end]), true), nil
	}
	return nil, nil
}

func (t *edgegram) next() (*jargon.Token, error) {
	// Emit from existing grams if available
	if tk, _ := t.nextFromGrams(); tk != nil {
		return tk, nil
	}

	for {
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
		vb := []byte(v)

		// Determine grapheme boundaries
		var (
			k   int
			pos []int
		)
		if isASCII(vb) {
			k = len(vb)
		} else {
			sg := graphemes.SegmentAll(vb)
			k = len(sg)
			pos = make([]int, k+1)
			for i := 1; i <= k; i++ {
				pos[i] = pos[i-1] + len(sg[i-1])
			}
		}

		if k == 0 {
			// No grams to emit; continue to next upstream token
			continue
		}

		// Build unified ranges: prefixes + windows
		n := t.size
		w := 0
		if n > 0 && k >= n {
			w = k - n + 1
		}
		ranges := make([]byteRange, 0, k+w)

		// Prefixes
		if pos == nil { // ASCII fast path
			for i := 1; i <= k; i++ {
				ranges = append(ranges, byteRange{start: 0, end: i})
			}
		} else {
			for i := 1; i <= k; i++ {
				ranges = append(ranges, byteRange{start: 0, end: pos[i]})
			}
		}

		// Windows of size n
		if n > 0 && k >= n {
			if pos == nil { // ASCII
				for i := n; i <= k; i++ {
					ranges = append(ranges, byteRange{start: i - n, end: i})
				}
			} else {
				for i := n; i <= k; i++ {
					ranges = append(ranges, byteRange{start: pos[i-n], end: pos[i]})
				}
			}
		}

		// Sort once by lexicographic order of underlying bytes
		sort.Slice(ranges, func(i, j int) bool {
			ri, rj := ranges[i], ranges[j]
			return bytes.Compare(vb[ri.start:ri.end], vb[rj.start:rj.end]) < 0
		})

		// Dedup adjacent equals
		grams := ranges[:0]
		for idx, r := range ranges {
			if idx == 0 {
				grams = append(grams, r)
				continue
			}
			prev := grams[len(grams)-1]
			if !bytes.Equal(vb[prev.start:prev.end], vb[r.start:r.end]) {
				grams = append(grams, r)
			}
		}

		// Initialize state for lazy emission
		t.vb = vb
		t.grams = grams
		t.cursor = 0

		if tk, _ := t.nextFromGrams(); tk != nil {
			return tk, nil
		}
		// If no grams (shouldn't happen if k>0), continue loop
	}
}

func isASCII(vb []byte) bool {
	for _, b := range vb {
		if b&0x80 != 0 {
			return false
		}
	}
	return true
}
