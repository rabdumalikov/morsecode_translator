package tests

import (
	"morse_converter/morse/encoding"
	"morse_converter/morse/mapping"
	"strings"
	"testing"
)

func TestEncodeDecoder(t *testing.T) {
	mapper, err := mapping.NewMapper("../../mapping/symbol2morse.json")
	if err != nil {
		t.Fatalf("Mapping creation failed with [%v]\n", err)
		return
	}
	encoder := encoding.NewEncoder(mapper)
	decoder := encoding.NewDecoder(mapper)

	tests := []struct {
		text     string
		expected string
	}{
		{"HELLO", "HELLO"},
		{"WORLD", "WORLD"},
		{"HELLO WORLD", "HELLO WORLD"},
		{"A A A A", "A A A A"},
		{"    Hello", "Hello"}, // separator last space
		{"Hello    ", "Hello"}, // separator first space
		{"    ", ""},           // no-separator
		{" Hello", "Hello"},    // separator last space
		{"Hello ", "Hello"},    // separator last space
		{"  ", ""},             // separator last space
		{"  547.33   _Dict._ v. _Gastroraphy_.[”]                   Removed.", "547.33 _Dict._ v. _Gastroraphy_.[”] Removed."},
		{"n  e", "n e"},
	}

	for _, test := range tests {
		result := encoder.Encode(test.text)
		result = decoder.Decode(result)
		result = strings.ToLower(result)
		if result != strings.ToLower(test.expected) {
			t.Errorf("Encode(%q) = [%q]; want [%q]", test.text, result, strings.ToLower(test.expected))
		}
	}
}
