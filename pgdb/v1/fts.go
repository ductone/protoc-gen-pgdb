package v1

import "github.com/doug-martin/goqu/v9/exp"

func FullTextSearchVectors(msg DBReflectMessage) string {
	//
	// TODO: add a "transcode" version for FTS data field
	// __transcode_version:1
	//
	// https://github.com/clipperhouse/jargon
	// desc := msg.DBReflect().Descriptor()
	// option 1:
	// do ngrams
	// do split
	// do stemming AND non-stemmed
	// "aaa:3 abb:3 "::tsvector
	// + READ side needs function
	//    websearch_to_tsquery
	//
	// option 2:
	//   ... do ngrams?
	//  use webserach()
	// to_tsvector?
	// (more or less what we do today)
	return ""
}

// expose function to do same stemming
// returns exp including websearch_to_tsquery()

func FullTextSerachQuery(input string) exp.Expression {

	return nil
}
