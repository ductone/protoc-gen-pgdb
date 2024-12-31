package v1

import (
	"bytes"
	"strings"

	"github.com/doug-martin/goqu/v9/exp"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/proto"
)

func MarshalNestedRecord(msg DBReflectMessage, opts ...RecordOptionsFunc) (exp.Record, error) {
	recs, err := msg.DBReflect().Record(opts...)
	if err != nil {
		return nil, err
	}
	return recs, nil
}

func cleanJSON(jsonBytes []byte) ([]byte, error) {
	// Quick guard: check if \u0000 exists in the input bytes
	if !bytes.Contains(jsonBytes, []byte("\\u0000")) {
		return jsonBytes, nil
	}

	// Decode and sanitize JSON
	var data interface{}
	if err := jsoniter.Unmarshal(jsonBytes, &data); err != nil {
		return nil, err
	}

	// Recursively sanitize
	data = sanitizeJSON(data)

	// Re-encode the JSON
	return jsoniter.Marshal(data)
}

func sanitizeJSON(input interface{}) interface{} {
	switch v := input.(type) {
	case map[string]interface{}:
		for key, value := range v {
			v[key] = sanitizeJSON(value)
		}
	case []interface{}:
		for i, value := range v {
			v[i] = sanitizeJSON(value)
		}
	case string:
		return strings.ReplaceAll(v, "\u0000", "")
	}
	return input
}

func MarshalProtoJSON(msg proto.Message) ([]byte, error) {
	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return cleanJSON(data)
}

func MarshalJSON(msg any) ([]byte, error) {
	data, err := jsoniter.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return cleanJSON(data)
}
