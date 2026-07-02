package main

import (
	"encoding/json"
	"os"
	"testing"
)

// TestParityGolden regenerates the parity dump in-process and compares it to
// the committed golden (pgdb/v1/testdata/parity_golden.json). The same golden
// is committed in the Rust protoc-gen-pgdb-rs repo, whose tests perform the
// same comparison — keeping both implementations database-compatible.
//
// Params holding JSON payloads (jsonb columns) are compared semantically:
// protojson output key order is randomized by design.
func TestParityGolden(t *testing.T) {
	goldenBytes, err := os.ReadFile("../../pgdb/v1/testdata/parity_golden.json")
	if err != nil {
		t.Fatalf("read golden: %v (regenerate with: go run ./cmd/pgdb-paritydump > pgdb/v1/testdata/parity_golden.json)", err)
	}

	var golden map[string]msgDump
	if err := json.Unmarshal(goldenBytes, &golden); err != nil {
		t.Fatalf("parse golden: %v", err)
	}

	got := map[string]msgDump{
		"animals.v1.Pet":         dumpMessage(fixturePet()),
		"animals.v1.ScalarValue": dumpMessage(fixtureScalarValue()),
		"zoo.v1.Shop":            dumpMessage(fixtureShop()),
	}

	for name, want := range golden {
		g, ok := got[name]
		if !ok {
			t.Errorf("%s: missing from regenerated dump", name)
			continue
		}
		if g.TableName != want.TableName {
			t.Errorf("%s: table_name: got %q want %q", name, g.TableName, want.TableName)
		}
		if len(g.CreateSchema) != len(want.CreateSchema) {
			t.Errorf("%s: schema count: got %d want %d", name, len(g.CreateSchema), len(want.CreateSchema))
		} else {
			for i := range g.CreateSchema {
				if g.CreateSchema[i] != want.CreateSchema[i] {
					t.Errorf("%s: schema[%d]:\n got: %s\nwant: %s", name, i, g.CreateSchema[i], want.CreateSchema[i])
				}
			}
		}
		if g.InsertSQL != want.InsertSQL {
			t.Errorf("%s: insert_sql:\n got: %s\nwant: %s", name, g.InsertSQL, want.InsertSQL)
		}
		compareParams(t, name+".insert", g.InsertParams, want.InsertParams)
		if g.DeleteSQL != want.DeleteSQL {
			t.Errorf("%s: delete_sql:\n got: %s\nwant: %s", name, g.DeleteSQL, want.DeleteSQL)
		}
		compareParams(t, name+".delete", g.DeleteParams, want.DeleteParams)
		if g.PKSK != want.PKSK {
			t.Errorf("%s: pksk: got %q want %q", name, g.PKSK, want.PKSK)
		}
	}
}

func compareParams(t *testing.T, prefix string, got, want []json.RawMessage) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("%s: param count: got %d want %d", prefix, len(got), len(want))
		return
	}
	for i := range got {
		if !paramsMatch(got[i], want[i]) {
			t.Errorf("%s: param %d: got %s want %s", prefix, i, got[i], want[i])
		}
	}
}

func paramsMatch(got, want json.RawMessage) bool {
	if string(got) == string(want) {
		return true
	}
	var g, w map[string]any
	if json.Unmarshal(got, &g) != nil || json.Unmarshal(want, &w) != nil {
		return false
	}
	if jsonEq(g, w) {
		return true
	}
	// jsonb payload strings: compare parsed
	gs, gok := g["s"].(string)
	ws, wok := w["s"].(string)
	if gok && wok && looksJSON(gs) && looksJSON(ws) {
		var gj, wj any
		if json.Unmarshal([]byte(gs), &gj) == nil && json.Unmarshal([]byte(ws), &wj) == nil {
			return jsonEq(gj, wj)
		}
	}
	return false
}

func looksJSON(s string) bool {
	return len(s) > 0 && (s[0] == '{' || s[0] == '[')
}

func jsonEq(a, b any) bool {
	switch av := a.(type) {
	case map[string]any:
		bv, ok := b.(map[string]any)
		if !ok || len(av) != len(bv) {
			return false
		}
		for k, v := range av {
			wv, ok := bv[k]
			if !ok || !jsonEq(v, wv) {
				return false
			}
		}
		return true
	case []any:
		bv, ok := b.([]any)
		if !ok || len(av) != len(bv) {
			return false
		}
		for i := range av {
			if !jsonEq(av[i], bv[i]) {
				return false
			}
		}
		return true
	case float64:
		bv, ok := b.(float64)
		return ok && av == bv
	default:
		return a == b
	}
}
