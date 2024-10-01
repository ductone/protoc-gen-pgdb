package swedish

import (
	"unicode/utf8"

	"github.com/kljensen/snowball/snowballword"
)

// Step 1 is the stemming of various endings found in
// R1 including "heterna", "ornas", and "andet".
func step1(w *snowballword.SnowballWord) bool {

	// Possible sufficies for this step, longest first.
	suffixes := []string{
		"heterna", "hetens", "anden", "heten", "heter", "arnas",
		"ernas", "ornas", "andes", "arens", "andet", "arna", "erna",
		"orna", "ande", "arne", "aste", "aren", "ades", "erns", "ade",
		"are", "ern", "ens", "het", "ast", "ad", "en", "ar", "er",
		"or", "as", "es", "at", "a", "e", "s",
	}

	// Using FirstSuffixIn since there are overlapping suffixes, where some might not be in the R1,
	// while another might. For example: "ärade"
	suffix := w.FirstSuffixIn(w.R1start, len(w.RS), suffixes...)
	suffixLength := utf8.RuneCountInString(suffix)

	// If it is not in R1, do nothing
	if suffix == "" || suffixLength > len(w.RS)-w.R1start {
		return false
	}

	if suffix == "s" {
		// Delete if preceded by a valid s-ending. Valid s-endings inlude the
		// following charaters: bcdfghjklmnoprtvy.
		//
		rsLen := len(w.RS)
		if rsLen >= 2 {
			switch w.RS[rsLen-2] {
			case 'b', 'c', 'd', 'f', 'g', 'h', 'j', 'k',
				'l', 'm', 'n', 'o', 'p', 'r', 't', 'v', 'y':
				w.RemoveLastNRunes(suffixLength)
				return true
			}
		}
		return false
	}
	// Remove the suffix
	w.RemoveLastNRunes(suffixLength)
	return true
}
