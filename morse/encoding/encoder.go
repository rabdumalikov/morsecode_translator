package encoding

import (
	"morse_converter/morse/mapping"
	"strings"
)

type Encoder struct {
	mapper mapping.Mapper
}

func NewEncoder(mapper mapping.Mapper) *Encoder {
	return &Encoder{mapper: mapper}
}

func (p *Encoder) Encode(text string) string {
	words := strings.Split(text, " ")
	var encodedWords []string

	for _, word := range words {
		var encodedWord []string

		for _, char := range word {
			upperChar := strings.ToUpper(string(char))
			value, ok := p.mapper.SymbolToMorse(upperChar)
			if ok {
				encodedWord = append(encodedWord, value)
			} else {
				encodedWord = append(encodedWord, string(char))
			}
		}
		encodedWords = append(encodedWords, strings.Join(encodedWord, " "))
	}

	return strings.Join(encodedWords, "/") + "//"
}
