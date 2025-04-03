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

func splitBySpace(text string) []string {

	var result []string
	for i, char := range text {

		if char == ' ' {
			// space is separator if there is a symbol on the left or on the right
			// checking left side

			doLeftSidePresent := (i-1 >= 0)
			nonSpaceOnTheLeft := false
			if i-1 >= 0 && text[i-1] != ' ' {
				nonSpaceOnTheLeft = true
			}

			// checking right side
			doRightSidePresent := (i+1 < len(text))
			nonSpaceOnTheRight := false
			if i+1 < len(text) && text[i+1] != ' ' {
				nonSpaceOnTheRight = true
			}

			isSpaceSeparator := (doLeftSidePresent && nonSpaceOnTheLeft) || (doLeftSidePresent && nonSpaceOnTheLeft && doRightSidePresent && nonSpaceOnTheRight)

			if isSpaceSeparator {
				result = append(result, string(""))
			} else {
				if len(result) == 0 {
					result = append(result, string(char))
				} else {
					result[len(result)-1] += string(char)
				}
			}

		} else {
			if len(result) == 0 {
				result = append(result, string(char))
			} else {
				result[len(result)-1] += string(char)
			}
		}

	}
	return result
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

		if len(encodedWord) != 0 {
			encodedWords = append(encodedWords, strings.Join(encodedWord, " "))
		}
	}

	return strings.Join(encodedWords, "/")
}
