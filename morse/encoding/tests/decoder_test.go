package tests

import (
	"morse_converter/morse/encoding"
	"morse_converter/morse/mapping"
	"testing"
)

func TestDecode(t *testing.T) {
	translator, err := mapping.NewTranslator("../../mapping/char2morse.json")
	if err != nil {
		t.Fatalf("Mapping creating failed with [%v]\n", err)
		return
	}
	decoder := encoding.NewDecoder(translator)

	tests := []struct {
		morseCode string
		expected  string
	}{
		{".... . .-.. .-.. ---", "HELLO"},
		{".... . .-.. .-.. ---/.-- --- .-. .-.. -..", "HELLO WORLD"},
	}

	for _, test := range tests {
		result, _ := decoder.DecodeLine(test.morseCode)
		if result != test.expected { // DecodeLine appends a newline
			t.Errorf("DecodeLine(%q) = [%q]; want [%q]", test.morseCode, result, test.expected+"\n")
		}
	}
}

func TestIncorrectInputDecode(t *testing.T) {
	translator, err := mapping.NewTranslator("../../mapping/char2morse.json")
	if err != nil {
		t.Fatalf("Mapping creating failed with [%v]\n", err)
		return
	}
	decoder := encoding.NewDecoder(translator)

	tests := []struct {
		morseCode string
		expected  string
	}{
		{".... . .-.. .-.. ---//", ""},
		{"//.... . .-.. .-.. ---", ""},
		{"////.... . .-.. .-.. ---", ""},
		{".... . .-.. .-.. ---//.-- --- .-. .-.. -..", ""},
	}

	for _, test := range tests {
		result, err := decoder.DecodeLine(test.morseCode)
		if err == nil || result != "" {
			t.Errorf("DecodeLine(%q) = err=[%v] result=[%q]; want non-empty error and empty result", test.morseCode, err, result)
		}
	}
}
