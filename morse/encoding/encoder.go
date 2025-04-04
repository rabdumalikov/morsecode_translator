package encoding

import (
	"morse_converter/morse/mapping"
	"strings"
)

type Encoder struct {
	mapper mapping.Mapper
}

const (
	AverageMorseCodeSize = 4 // symbols
)

func NewEncoder(mapper mapping.Mapper) *Encoder {
	return &Encoder{mapper: mapper}
}

func (p *Encoder) Encode(text string) string {
	var sb strings.Builder

	// Pre-allocate buffer
	sb.Grow(len(text) * AverageMorseCodeSize)

	needsCharSeparator := false

	for _, r := range text {
		if r == ' ' {
			// If the previous character was not a space (i.e. we finished a word)
			if needsCharSeparator {
				sb.WriteByte('/')
				needsCharSeparator = false
			}
			// Ignore consecutive spaces, only write one '/'
		} else {
			if needsCharSeparator {
				sb.WriteByte(' ') // Write character separator before the next Morse symbol
			}

			upperCharStr := strings.ToUpper(string(r))
			morseCode, ok := p.mapper.SymbolToMorse(upperCharStr)

			if ok {
				sb.WriteString(morseCode)
			} else {
				sb.WriteRune(r)
			}

			needsCharSeparator = true
		}
	}

	encodedText := strings.TrimSuffix(sb.String(), "/")

	return encodedText
}
