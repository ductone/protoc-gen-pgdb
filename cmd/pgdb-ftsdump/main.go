package main

// pgdb-ftsdump dumps Go FTS tsvector strings for a corpus of SearchContent
// docs (JSON lines on stdin). Used to regenerate the cross-language FTS
// parity fixtures shared with the Rust repo:
//   go run ./cmd/pgdb-ftsdump < pgdb/v1/testdata/fts_corpus.jsonl > pgdb/v1/testdata/fts_expected.jsonl
//
// Dumps Go FTS tsvector strings for a corpus, one JSON object per line:
// input: {"type": 1, "weight": 2, "value": "text"}  (a single-doc set)
// output: {"tsvector": "..."} — the raw string passed to ?::tsvector.
import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

type inputDoc struct {
	Type   int32  `json:"type"`
	Weight int32  `json:"weight"`
	Value  string `json:"value"`
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		var docs []inputDoc
		if err := json.Unmarshal(line, &docs); err != nil {
			panic(err)
		}
		sc2 := make([]*pgdb_v1.SearchContent, 0, len(docs))
		for _, d := range docs {
			sc2 = append(sc2, &pgdb_v1.SearchContent{
				Type:   pgdb_v1.FieldOptions_FullTextType(d.Type),
				Weight: pgdb_v1.FieldOptions_FullTextWeight(d.Weight),
				Value:  d.Value,
			})
		}
		e := pgdb_v1.FullTextSearchVectors(sc2)
		// render via goqu to extract the param
		ds := goqu.Dialect("postgres").From("t").Where(exp.NewLiteralExpression("?", e)).Prepared(true)
		sqlText, params, err := ds.ToSQL()
		if err != nil {
			panic(err)
		}
		v := ""
		empty := false
		if len(params) > 0 {
			v, _ = params[0].(string)
		} else {
			empty = true
		}
		_ = sqlText
		b, _ := json.Marshal(map[string]interface{}{"tsvector": v, "empty": empty})
		fmt.Fprintf(out, "%s\n", b)
	}
}
