package v1

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestSearchGithub(t *testing.T) {
	spew.Dump(FullTextSearchVectors([]SearchContent{
		{
			Type:   FieldOptions_FULL_TEXT_TYPE_EDGE_NGRAM,
			Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
			Value:  "github.com/pquerna/reponame",
		},
	}))
	spew.Dump(FullTextSearchVectors([]SearchContent{
		{
			Type:   FieldOptions_FULL_TEXT_TYPE_EDGE_NGRAM,
			Weight: FieldOptions_FULL_TEXT_WEIGHT_HIGH,
			Value:  "role/biquery.user",
		},
	}))
}
