package encoding

import (
	"morse_converter/morse/mapping"
	"strings"
)

type Decoder struct {
	mapper mapping.Mapper
}

func NewDecoder(mapper mapping.Mapper) *Decoder {
	return &Decoder{mapper: mapper}
}

func (p *Decoder) Decode(morseCode string) string {
	words := strings.Split(morseCode, "/")
	var decodedWords []string

	for _, word := range words {
		if word == "" {
			continue // Skip empty words
		}

		var decodedChars []string
		for _, code := range strings.Split(strings.TrimSpace(word), " ") {
			if code == "" {
				continue
			}

			value, ok := p.mapper.MorseToSymbol(code)
			if ok {
				decodedChars = append(decodedChars, value)
			} else {
				decodedChars = append(decodedChars, code) // Keep unknown symbols
			}
		}
		decodedWords = append(decodedWords, strings.Join(decodedChars, ""))
	}

	return strings.Join(decodedWords, " ") + "\n"
}
