package encoding

import (
	"fmt"
	"morse_converter/morse/mapping"
	"strings"
)

type Encoder struct {
	translator mapping.Translator
}

const (
	AverageMorseCodeSize = 4 // Symbols
)

func NewEncoder(translator mapping.Translator) *Encoder {
	return &Encoder{translator: translator}
}

// EncodeLine encodes a single text line into morse code.
// It does not support multiple sentences or text new line symbols (i.e. "\n").
// Characters that have no Morse code representation would be copied as-is.
func (p *Encoder) EncodeLine(text string) (string, error) {
	// Validate input
	if strings.Contains(text, string(mapping.TextNewLine)) {
		return "", fmt.Errorf("input contains unsupported new line '%c' separator", mapping.TextNewLine)
	}

	var sb strings.Builder

	// Pre-allocate buffer
	sb.Grow(len(text) * AverageMorseCodeSize)

	needsCharSeparator := false

	for _, char := range text {
		if char == rune(mapping.TextWordSeparator) {
			// If the previous character was not a space (i.e. we finished a word)
			if needsCharSeparator {
				sb.WriteByte(mapping.MorseWordSeparator)
				needsCharSeparator = false
			}
			// Ignore consecutive spaces, only write one '/'
		} else {
			if needsCharSeparator {
				sb.WriteByte(mapping.TextWordSeparator) // Write character separator before the next Morse code
			}

			upperCharStr := strings.ToUpper(string(char))
			morseCode, ok := p.translator.CharToMorse(upperCharStr)

			if ok {
				sb.WriteString(morseCode)
			} else {
				// Copy char that has no morse code representation as-is
				sb.WriteRune(char)
			}

			needsCharSeparator = true
		}
	}

	encodedText := strings.TrimSuffix(sb.String(), string(mapping.MorseWordSeparator))

	return encodedText, nil
}
