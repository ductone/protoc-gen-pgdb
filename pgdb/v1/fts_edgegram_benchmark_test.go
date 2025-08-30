package v1

import (
	"fmt"
	"strings"
	"testing"

	"github.com/clipperhouse/jargon"
)

// singleWordStream returns a TokenStream that yields exactly one word token, then EOF.
func singleWordStream(s string) *jargon.TokenStream {
	emitted := false
	word := jargon.NewToken(s, true)
	return jargon.NewTokenStream(func() (*jargon.Token, error) {
		if emitted {
			return nil, nil
		}
		emitted = true
		return word, nil
	})
}

func benchmarkInputs() []struct{ name, s string } {
	return []struct{ name, s string }{
		{name: "SmallASCII_git", s: "git"},
		{name: "SmallASCII_hub", s: "hub"},
		{name: "MediumASCII_github", s: "github"},
		{name: "Sentence", s: "super sweet"},
		{name: "Punctuated", s: "go, go-go!"},
		{name: "UnicodeCombining", s: "cafe\u0301"},
		{name: "EmojiZWJ", s: "👨‍👩‍👧‍👦"},
		{name: "MixedScripts", s: "東京toronto"},
		{name: "LongASCII_1KiB", s: strings.Repeat("a", 1024)},
		{name: "LongASCII_4KiB", s: strings.Repeat("github", 512)},
		{name: "Paragraph_English_2KiB", s: strings.Repeat("The quick brown fox jumps over the lazy dog. ", 40)},
	}
}

func nValues() []int { return []int{1, 2, 3, 5, 8, 16} }

// BenchmarkEdgegramStream_Pipeline measures the end-to-end pipeline cost including tokenization.
func BenchmarkEdgegramStream_Pipeline(b *testing.B) {
	b.ReportAllocs()
	for _, c := range benchmarkInputs() {
		for _, n := range nValues() {
			b.Run(fmt.Sprintf("%s/n=%d", c.name, n), func(b *testing.B) {
				filter := edgegramStream(n)
				b.ResetTimer()
				var sink int
				for i := 0; i < b.N; i++ {
					tokens, err := jargon.TokenizeString(c.s).Filter(filter).Words().ToSlice()
					if err != nil {
						b.Fatal(err)
					}
					sink += len(tokens)
				}
				if sink == -1 {
					b.Fatal("unreachable")
				}
			})
		}
	}
}

// BenchmarkEdgegramStream_Isolated measures the cost of the edgegram filter independent of tokenization.
func BenchmarkEdgegramStream_Isolated(b *testing.B) {
	b.ReportAllocs()
	for _, c := range benchmarkInputs() {
		for _, n := range nValues() {
			b.Run(fmt.Sprintf("%s/n=%d", c.name, n), func(b *testing.B) {
				filter := edgegramStream(n)
				// Prebuild the input token outside the timed region.
				b.ResetTimer()
				var sink int
				for i := 0; i < b.N; i++ {
					incoming := singleWordStream(c.s)
					outgoing := filter(incoming)
					for {
						tk, err := outgoing.Next()
						if err != nil {
							b.Fatal(err)
						}
						if tk == nil {
							break
						}
						// Accumulate to avoid dead-code elimination.
						sink += len(tk.String())
					}
				}
				if sink == -1 {
					b.Fatal("unreachable")
				}
			})
		}
	}
}
