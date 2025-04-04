package tests

import (
	"morse_converter/morse/encoding"
	"morse_converter/morse/mapping"
	"testing"
)

func TestEncode(t *testing.T) {
	mapper, err := mapping.NewMapper("../../mapping/symbol2morse.json")
	if err != nil {
		t.Fatalf("Mapping creation failed with [%v]\n", err)
		return
	}
	encoder := encoding.NewEncoder(mapper)

	tests := []struct {
		text     string
		expected string
	}{
		{"HELLO", ".... . .-.. .-.. ---"},
		{"WORLD", ".-- --- .-. .-.. -.."},
		{"HELLO WORLD", ".... . .-.. .-.. ---/.-- --- .-. .-.. -.."},
		{"A A A A", ".-/.-/.-/.-"},
		{"    Hello", ".... . .-.. .-.. ---"},
		{"Hello    ", ".... . .-.. .-.. ---"},
		{"    ", ""},
		{" Hello", ".... . .-.. .-.. ---"},
		{"Hello ", ".... . .-.. .-.. ---"},
		{"\302\240", "\302\240"},
	}

	for _, test := range tests {
		result := encoder.Encode(test.text)
		if result != test.expected {
			t.Errorf("Encode(%q) = [%q]; want [%q]", test.text, result, test.expected)
		}
	}
}
