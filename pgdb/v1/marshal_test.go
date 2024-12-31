package v1

import (
	reflect "reflect"
	"strings"
	"testing"

	pstruct "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/protobuf/encoding/protojson"

	jsoniter "github.com/json-iterator/go"
)

func TestCleanJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "No changes needed",
			input:    `{"key": "value"}`,
			expected: `{"key": "value"}`,
		},
		{
			name:     "Remove single \u0000",
			input:    `{"key": "value\u0000"}`,
			expected: `{"key":"value"}`,
		},
		{
			name:     "Nested object with \u0000",
			input:    `{"outer": {"inner": "data\u0000"}}`,
			expected: `{"outer":{"inner":"data"}}`,
		},
		{
			name:     "Array with \u0000",
			input:    `{"list": ["item1", "item2\u0000"]}`,
			expected: `{"list":["item1","item2"]}`,
		},
		{
			name:     "Multiple \u0000",
			input:    `{"key": "val\u0000ue\u0000"}`,
			expected: `{"key":"value"}`,
		},
		{
			name:     "Colon delimited string",
			input:    `{"key": "111":222:3333:4444"}`,
			expected: `{"key": "111":222:3333:4444"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := cleanJSON([]byte(test.input))
			if (err != nil) != test.hasError {
				t.Errorf("Unexpected error status. Got: %v, Expected error: %v", err, test.hasError)
				return
			}
			if !test.hasError && string(result) != test.expected {
				t.Errorf("Unexpected result. Got: %s, Expected: %s", result, test.expected)
			}
		})
	}
}

func TestMarshalProtoJSON(t *testing.T) {
	msg := &pstruct.Struct{
		Fields: map[string]*pstruct.Value{
			"PK:SK": {
				Kind: &pstruct.Value_StringValue{
					StringValue: "uid_kdjkdjjdki:uid_ksdjfksdjfkjs:sdflsdfjkl\u0000",
				},
			},
			"SomeString": {
				Kind: &pstruct.Value_StringValue{
					StringValue: "\u0000uid_kdjkdjjdki:uid\u0000_ksdjfksdjfkjs:sdfl\u0000sdfjkl\u0000",
				},
			},
		},
	}

	msg2 := &pstruct.Struct{
		Fields: map[string]*pstruct.Value{
			"PK:SK": {
				Kind: &pstruct.Value_StringValue{
					StringValue: "uid_kdjkdjjdki:uid_ksdjfksdjfkjs:sdflsdfjkl",
				},
			},
			"SomeString": {
				Kind: &pstruct.Value_StringValue{
					StringValue: "uid_kdjkdjjdki:uid_ksdjfksdjfkjs:sdflsdfjkl",
				},
			},
		},
	}
	msg_expected := &pstruct.Struct{
		Fields: map[string]*pstruct.Value{
			"PK:SK": {
				Kind: &pstruct.Value_StringValue{
					StringValue: "uid_kdjkdjjdki:uid_ksdjfksdjfkjs:sdflsdfjkl",
				},
			},
			"SomeString": {
				Kind: &pstruct.Value_StringValue{
					StringValue: "uid_kdjkdjjdki:uid_ksdjfksdjfkjs:sdflsdfjkl",
				},
			},
		},
	}

	msg_expected_bytes, err := protojson.Marshal(msg_expected)
	if err != nil {
		t.Errorf("Failed to marshal expected message: %v", err)
	}

	expected := string(msg_expected_bytes)
	expected_no_space := strings.ReplaceAll(expected, " ", "")

	tests := []struct {
		name     string
		input    *pstruct.Struct
		expected string
	}{
		{
			name:     "With null bytes",
			input:    msg,
			expected: expected_no_space,
		},
		{
			name:     "Without null bytes",
			input:    msg2,
			expected: expected,
		},
	}

	result, err := MarshalProtoJSON(msg)
	if err != nil {
		t.Fatalf("Expected an error when interpreting binary proto data as JSON, got none.\nResult: %s", string(result))
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := MarshalProtoJSON(test.input)
			if err != nil {
				t.Fatalf("Expected an error when interpreting binary proto data as JSON, got none.\nResult: %s", string(result))
			}

			if string(result) != test.expected {
				t.Errorf("Unexpected result.\nGot:      %s\nExpected: %s", result, test.expected)
			}
		})
	}
}

func TestMarshalJSON(t *testing.T) {
	msg := map[string]interface{}{
		"PK:SK":      "uid_kdjkdjjdki:uid_ksdjfksdjfkjs:sdflsdfjkl",
		"SomeString": "\u0000String!!!! \u0000with\u0000",
	}

	msg2 := map[string]interface{}{
		"PK:SK":      "uid_kdjkdjjdki:uid_ksdjfksdjfkjs:sdflsdfjkl",
		"SomeString": "String!!!! with",
	}

	tests := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name:     "With null bytes",
			input:    msg,
			expected: msg2,
		},
		{
			name:     "Without null bytes",
			input:    msg2,
			expected: msg2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := MarshalJSON(test.input)
			if err != nil {
				t.Fatalf("Expected an error when interpreting binary proto data as JSON, got none.\nResult: %s", string(result))
			}

			var got map[string]interface{}
			if err := jsoniter.Unmarshal(result, &got); err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			if !reflect.DeepEqual(got, test.expected) {
				t.Errorf("Unexpected result.\nGot:      %v\nExpected: %v", got, test.expected)
			}
		})
	}
}
