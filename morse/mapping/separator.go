package mapping

const (
	MorseWordSeparator byte   = '/'
	MorseCharSeparator byte   = ' '
	MorseNewLine       string = "//"

	TextWordSeparator byte   = ' '
	TextCharSeparator string = ""
	TextNewLine       byte   = '\n'
)

type SeparatorType int

const (
	WordSeparator SeparatorType = iota
	CharSeparator
	NewLineSeparator
)

type Separator interface {
	ToString() string
}

type TextSeparator interface {
	Separator
	ToMorse() MorseSeparator
}

type MorseSeparator interface {
	Separator
	ToText() TextSeparator
}

type textSeparator struct {
	separator SeparatorType
}

func NewTextSeparator(separator SeparatorType) TextSeparator {
	return &textSeparator{separator: separator}
}

func (ts *textSeparator) ToString() string {
	switch ts.separator {
	case WordSeparator:
		return string(TextWordSeparator)
	case CharSeparator:
		return string(TextCharSeparator)
	case NewLineSeparator:
		return string(TextNewLine)
	default: // This will never be the case
		return ""
	}
}

func (ts *textSeparator) ToMorse() MorseSeparator {
	return &morseSeparator{separator: ts.separator}
}

type morseSeparator struct {
	separator SeparatorType
}

func NewMorseSeparator(separator SeparatorType) MorseSeparator {
	return &morseSeparator{separator: separator}
}

func (ms *morseSeparator) ToString() string {
	switch ms.separator {
	case WordSeparator:
		return string(MorseWordSeparator)
	case CharSeparator:
		return string(MorseCharSeparator)
	case NewLineSeparator:
		return MorseNewLine
	default: // This will never be the case
		return ""
	}
}

func (ms *morseSeparator) ToText() TextSeparator {
	return &textSeparator{separator: ms.separator}
}
