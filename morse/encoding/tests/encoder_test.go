package tests

import (
	"morse_converter/morse/encoding"
	"morse_converter/morse/mapping"
	"testing"
)

func TestEncode(t *testing.T) {
	translator, err := mapping.NewTranslator("../../mapping/char2morse.json")
	if err != nil {
		t.Fatalf("Mapping creation failed with [%v]\n", err)
		return
	}
	encoder := encoding.NewEncoder(translator)

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
		result, _ := encoder.EncodeLine(test.text)
		if result != test.expected {
			t.Errorf("EncodeLine(%q) = [%q]; want [%q]", test.text, result, test.expected)
		}
	}
}

func TestIncorrectInputEncode(t *testing.T) {
	translator, err := mapping.NewTranslator("../../mapping/char2morse.json")
	if err != nil {
		t.Fatalf("Mapping creation failed with [%v]\n", err)
		return
	}
	encoder := encoding.NewEncoder(translator)

	tests := []struct {
		text     string
		expected string
	}{
		{"HELLO\n", ".... . .-.. .-.. ---"},
		{"WORLD\n", ".-- --- .-. .-.. -.."},
		{"HELLO\nWORLD", ".... . .-.. .-.. ---/.-- --- .-. .-.. -.."},
		{"A'nA\nA\nA", ".-/.-/.-/.-"},
		{"\n\n\nHello", ".... . .-.. .-.. ---"},
		{"Hello\r\n", ".... . .-.. .-.. ---"},
	}

	for _, test := range tests {
		result, err := encoder.EncodeLine(test.text)
		if err == nil || result != "" {
			t.Errorf("EncodeLine(%q) = err=[%v] result=[%q]; want non-empty error and empty result", test.text, err, result)
		}
	}
}
