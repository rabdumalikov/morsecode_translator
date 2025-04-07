package encoding

import (
	"fmt"
	"morse_converter/morse/mapping"
	"strings"
)

type Decoder struct {
	translator mapping.Translator
}

func NewDecoder(translator mapping.Translator) *Decoder {
	return &Decoder{translator: translator}
}

// DecodeLine decodes a single Morse code line into text.
// It does not support multiple sentences or Morse new line symbols (i.e. "//").
// Characters that have no Morse code representation would be copied as-is.
func (p *Decoder) DecodeLine(morseCode string) (string, error) {
	// Validate input
	if strings.Contains(morseCode, mapping.MorseNewLine) {
		return "", fmt.Errorf("input contains unsupported new line '%s' separator", mapping.MorseNewLine)
	}

	words := strings.Split(morseCode, string(mapping.MorseWordSeparator))
	var decodedWords []string

	for _, word := range words {
		if word == "" {
			continue // Skip empty words
		}

		var decodedChars []string
		for _, code := range strings.Split(word, string(mapping.MorseCharSeparator)) {
			if code == "" {
				continue
			}

			value, ok := p.translator.MorseToChar(code)
			if ok {
				decodedChars = append(decodedChars, value)
			} else {
				decodedChars = append(decodedChars, code) // Keep unknown symbols
			}
		}
		decodedWords = append(decodedWords, strings.Join(decodedChars, ""))
	}

	return strings.Join(decodedWords, string(mapping.TextWordSeparator)), nil
}
