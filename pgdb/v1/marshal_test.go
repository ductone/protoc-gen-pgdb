package v1

import "testing"

// Test cases for CleanJSON
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
