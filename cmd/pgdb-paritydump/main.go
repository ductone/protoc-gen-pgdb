// pgdb-paritydump emits a deterministic JSON dump of the DB-facing outputs
// (DDL, insert/delete SQL and params, PKSK) for a set of fixture messages,
// under DialectV17.
//
// The output is committed as a golden fixture in BOTH the Go repo and the
// Rust protoc-gen-pgdb-rs repo: each side's tests regenerate their dump and
// compare against it, guaranteeing the two implementations stay
// database-compatible.
//
// Regenerate with: go run ./cmd/pgdb-paritydump > pgdb/v1/testdata/parity_golden.json
package main

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"

	animals_v1 "github.com/ductone/protoc-gen-pgdb/example/models/animals/v1"
	zoo_v1 "github.com/ductone/protoc-gen-pgdb/example/models/zoo/v1"
	pgdb_v1 "github.com/ductone/protoc-gen-pgdb/pgdb/v1"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type msgDump struct {
	TableName    string            `json:"table_name"`
	CreateSchema []string          `json:"create_schema"`
	InsertSQL    string            `json:"insert_sql"`
	InsertParams []json.RawMessage `json:"insert_params"`
	DeleteSQL    string            `json:"delete_sql"`
	DeleteParams []json.RawMessage `json:"delete_params"`
	PKSK         string            `json:"pksk"`
}

func encParam(v any) json.RawMessage {
	enc := func(m map[string]any) json.RawMessage {
		b, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}
		return b
	}
	if v == nil {
		return enc(map[string]any{"null": true})
	}
	switch t := v.(type) {
	case string:
		return enc(map[string]any{"s": t})
	case *string:
		if t == nil {
			return enc(map[string]any{"null": true})
		}
		return enc(map[string]any{"s": *t})
	case []byte:
		return enc(map[string]any{"x": hex.EncodeToString(t)})
	case bool:
		return enc(map[string]any{"bool": t})
	case int, int32, int64:
		return enc(map[string]any{"i": t})
	case uint, uint32, uint64:
		return enc(map[string]any{"u": t})
	case float32, float64:
		return enc(map[string]any{"f": t})
	case time.Time:
		return enc(map[string]any{"t": t.UTC().Format(time.RFC3339Nano)})
	case *time.Time:
		if t == nil {
			return enc(map[string]any{"null": true})
		}
		return enc(map[string]any{"t": t.UTC().Format(time.RFC3339Nano)})
	case pgtype.Interval:
		return enc(map[string]any{"interval_us": t.Microseconds, "valid": t.Valid})
	case *pgtype.Interval:
		return enc(map[string]any{"interval_us": t.Microseconds, "valid": t.Valid})
	case driver.Valuer:
		dv, err := t.Value()
		if err != nil {
			panic(err)
		}
		if dv == nil {
			return enc(map[string]any{"null": true})
		}
		return encParam(dv)
	default:
		return enc(map[string]any{"unknown": fmt.Sprintf("%T", v)})
	}
}

func encParams(params []any) []json.RawMessage {
	rv := make([]json.RawMessage, 0, len(params))
	for _, p := range params {
		rv = append(rv, encParam(p))
	}
	return rv
}

// insertParamIndex finds the bind-parameter index (0-based) of a column in a
// single-row INSERT by locating the column position and reading the matching
// VALUES item (columns rendered as NULL literals consume no parameter).
func insertParamIndex(insertSQL, col string) int {
	start := strings.Index(insertSQL, "(")
	end := strings.Index(insertSQL, ") VALUES ")
	if start < 0 || end < 0 {
		return -1
	}
	cols := strings.Split(insertSQL[start+1:end], ", ")
	colIdx := -1
	for i, c := range cols {
		if c == `"`+col+`"` {
			colIdx = i
			break
		}
	}
	if colIdx < 0 {
		return -1
	}
	valuesPart := insertSQL[end+len(") VALUES "):]
	if onConflict := strings.Index(valuesPart, " ON CONFLICT"); onConflict >= 0 {
		valuesPart = valuesPart[:onConflict]
	}
	valuesPart = strings.TrimPrefix(valuesPart, "(")
	valuesPart = strings.TrimSuffix(valuesPart, ")")
	items := strings.Split(valuesPart, ", ")
	if colIdx >= len(items) {
		return -1
	}
	item := items[colIdx]
	dollar := strings.Index(item, "$")
	if dollar < 0 {
		return -1
	}
	n := 0
	for _, r := range item[dollar+1:] {
		if r < '0' || r > '9' {
			break
		}
		n = n*10 + int(r-'0')
	}
	if n == 0 {
		return -1
	}
	return n - 1
}

func dumpMessage(msg pgdb_v1.DBReflectMessage) msgDump {
	dialect := pgdb_v1.DialectV17
	dbr := msg.DBReflect(dialect)
	desc := dbr.Descriptor()

	schema, err := pgdb_v1.CreateSchema(msg, dialect)
	if err != nil {
		panic(err)
	}

	insertSQL, insertParams, err := pgdb_v1.Insert(msg, dialect)
	if err != nil {
		panic(err)
	}

	deleteSQL, deleteParams, err := pgdb_v1.Delete(msg, dialect)
	if err != nil {
		panic(err)
	}

	pksk := ""
	if pk, ok := dbr.(pgdb_v1.PrimaryKeyer); ok {
		pksk = pk.PKSK()
	}

	// Replace pb$pb_data with a deterministic marshal so the golden is stable
	// across runs and languages (proto.Marshal map ordering is randomized).
	if pm, ok := msg.(proto.Message); ok {
		if idx := insertParamIndex(insertSQL, "pb$pb_data"); idx >= 0 && idx < len(insertParams) {
			det, err := proto.MarshalOptions{Deterministic: true}.Marshal(pm)
			if err != nil {
				panic(err)
			}
			insertParams[idx] = det
		}
	}

	return msgDump{
		TableName:    desc.TableName(),
		CreateSchema: schema,
		InsertSQL:    insertSQL,
		InsertParams: encParams(insertParams),
		DeleteSQL:    deleteSQL,
		DeleteParams: encParams(deleteParams),
		PKSK:         pksk,
	}
}

func fixturePet() *animals_v1.Pet {
	profile, err := structpb.NewStruct(map[string]any{
		"color":      "brown",
		"age":        3.5,
		"vaccinated": true,
	})
	if err != nil {
		panic(err)
	}
	extra, err := structpb.NewStruct(map[string]any{"a": 1.5})
	if err != nil {
		panic(err)
	}
	b := animals_v1.Pet_builder{
		TenantId:    "t1",
		Id:          "pet_1",
		CreatedAt:   timestamppb.New(time.Unix(1672628645, 0).UTC()),
		UpdatedAt:   timestamppb.New(time.Unix(1672628646, 0).UTC()),
		DisplayName: "Fluffy McFlufferson",
		Description: "A very good dog. Best boy!",
		SystemBuiltin: true,
		Elapsed:     durationpb.New(90*time.Second + 500*time.Millisecond),
		Profile:     profile,
		Cuteness:    0.9,
		Price:       199.99,
		VeryLongNaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaame: true,
		ExtraProfiles: []*structpb.Struct{extra},
		FieldWithV17CollationOnly: "2K5uidLikeValue",
	}
	return b.Build()
}

func fixtureScalarValue() *animals_v1.ScalarValue {
	b := animals_v1.ScalarValue_builder{
		TenantId: "t1",
		Id:       "sv_1",
		Double:   1.5,
		Float:    2.25,
		Int32:    -32,
		Int64:    -64,
		Uint32:   32,
		Uint64:   18446744073709551615,
		Sint32:   -320,
		Sint64:   -640,
		Fixed32:  1132,
		Fixed64:  1164,
		Sfixed32: -1132,
		Sfixed64: -1164,
		Bool:     true,
		String:   "Hello postgres World",
		Bytes:    []byte{0xde, 0xad, 0xbe, 0xef},

		RepeatedDouble:   []float64{1.5, -2.5},
		RepeatedFloat:    []float32{0.5},
		RepeatedInt32:    []int32{1, -2, 3},
		RepeatedInt64:    []int64{4, -5},
		RepeatedUint32:   []uint32{6},
		RepeatedUint64:   []uint64{7, 8},
		RepeatedSint32:   []int32{-9},
		RepeatedSint64:   []int64{10},
		RepeatedFixed32:  []uint32{11},
		RepeatedFixed64:  []uint64{12},
		RepeatedSfixed32: []int32{-13},
		RepeatedSfixed64: []int64{14},
		RepeatedBool:     []bool{true, false},
		RepeatedString:   []string{"alpha", "beta gamma", "日本語"},
		RepeatedBytes:    [][]byte{{0x01}, {0x02, 0x03}},
		RepeatedEnum:     []animals_v1.FurType{animals_v1.FurType_FUR_TYPE_LOTS},

		StringMap: map[string]string{"k1": "v1"},

		CreatedAt: timestamppb.New(time.Unix(1672628645, 0).UTC()),
	}
	return b.Build()
}

func fixtureShop() *zoo_v1.Shop {
	b := zoo_v1.Shop_builder{
		TenantId:  "t1",
		Id:        "shop_1",
		CreatedAt: timestamppb.New(time.Unix(1672628645, 0).UTC()),
		Anything:  fixtureScalarValue(),
		Fur:       animals_v1.FurType_FUR_TYPE_LOTS,
		Mgr:       zoo_v1.Shop_Manager_builder{Id: 7}.Build(),
	}
	return b.Build()
}

func main() {
	out := map[string]msgDump{
		"animals.v1.Pet":         dumpMessage(fixturePet()),
		"animals.v1.ScalarValue": dumpMessage(fixtureScalarValue()),
		"zoo.v1.Shop":            dumpMessage(fixtureShop()),
	}

	// deterministic key order
	keys := make([]string, 0, len(out))
	for k := range out {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	ordered := make(map[string]msgDump, len(out))
	for _, k := range keys {
		ordered[k] = out[k]
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(ordered); err != nil {
		panic(err)
	}
}
