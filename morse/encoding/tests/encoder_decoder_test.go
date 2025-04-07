package tests

import (
	"morse_converter/morse/encoding"
	"morse_converter/morse/mapping"
	"strings"
	"testing"
)

func TestEncodeDecoder(t *testing.T) {
	translator, err := mapping.NewTranslator("../../mapping/char2morse.json")
	if err != nil {
		t.Fatalf("Mapping creation failed with [%v]\n", err)
		return
	}
	encoder := encoding.NewEncoder(translator)
	decoder := encoding.NewDecoder(translator)

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
		result, _ := encoder.EncodeLine(test.text)
		result, _ = decoder.DecodeLine(result)
		result = strings.ToLower(result)
		if result != strings.ToLower(test.expected) {
			t.Errorf("EncodeLine(%q) = [%q]; want [%q]", test.text, result, strings.ToLower(test.expected))
		}
	}
}
