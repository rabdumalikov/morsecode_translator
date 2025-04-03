package tests

import (
	"morse_converter/morse/encoding"
	"morse_converter/morse/mapping"
	"testing"
)

func TestDecode(t *testing.T) {
	mapper, err := mapping.NewMapper("../../mapping/symbol2morse.json")
	if err != nil {
		t.Fatalf("Mapping creating failed with [%v]\n", err)
		return
	}
	decoder := encoding.NewDecoder(mapper)

	tests := []struct {
		morseCode string
		expected  string
	}{
		{".... . .-.. .-.. ---", "HELLO"},
		{".... . .-.. .-.. ---/.-- --- .-. .-.. -..", "HELLO WORLD"},
	}

	for _, test := range tests {
		result := decoder.Decode(test.morseCode)
		if result != test.expected { // Decode appends a newline
			t.Errorf("Decode(%q) = [%q]; want [%q]", test.morseCode, result, test.expected+"\n")
		}
	}
}
