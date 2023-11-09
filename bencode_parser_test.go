package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestParseBencoding(t *testing.T) {
	testCases := []struct {
		input    string
		expected interface{}
		err      error
	}{
		{"", nil, errors.New("Empty input data")},
		{"i123e", 123, nil},
		{"5:hello", "hello", nil},
		{"l3:onei2ee", []interface{}{"one", 2}, nil},
		{"d3:key5:valuee", map[string]interface{}{"key": "value"}, nil},
		{"invalid_input", 0, errors.New("Invalid integer format")},
	}

	for _, tc := range testCases {
		result, err := parseBencoding(tc.input)

		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("Input: %s, Expected: %v, Got: %v", tc.input, tc.expected, result)
		}

		if (err == nil && tc.err != nil) || (err != nil && tc.err == nil) || (err != nil && tc.err != nil && err.Error() != tc.err.Error()) {
			t.Errorf("Input: %s, Expected Error: %v, Got Error: %v", tc.input, tc.err, err)
		}
	}
}

// Add more tests for other functions if needed
