package mapping

import (
	"encoding/json"
	"maps"
	"os"
)

type Mapper interface {
	SymbolToMorse(symbol string) (string, bool)
	MorseToSymbol(morse string) (string, bool)
	MorseWordSeparator() string
	MorseCodeSeparator() string
	MorseNewLineSeparator() string
	ToTextSeparator(morseSeparator string) string
	ToMorseSeparator(textSeparator string) string
}

type Mapping struct {
	mapping              rawMapping
	symbolToMorseMapping map[string]string
	morseToSymbolMapping map[string]string
}

func (p *Mapping) SymbolToMorse(symbol string) (string, bool) {
	if p.symbolToMorseMapping == nil {
		p.symbolToMorseMapping = createSymbolToMorseTable(p.mapping)
	}

	value, ok := p.symbolToMorseMapping[symbol]
	return value, ok
}

func (p *Mapping) MorseToSymbol(morse string) (string, bool) {
	if p.morseToSymbolMapping == nil {
		p.morseToSymbolMapping = createMorseToSymbolTable(p.mapping)
	}

	value, ok := p.morseToSymbolMapping[morse]
	return value, ok
}

func (p *Mapping) MorseWordSeparator() string {
	return "/"
}

func (p *Mapping) MorseCodeSeparator() string {
	return " "
}

func (p *Mapping) MorseNewLineSeparator() string {
	return "//"
}

func (p *Mapping) ToTextSeparator(morseSeparator string) string {
	switch morseSeparator {
	case "/":
		return " "
	case "//":
		return "\n"
	default:
		return ""
	}
}

func (p *Mapping) ToMorseSeparator(textSeparator string) string {
	switch textSeparator {
	case " ":
		return "/"
	case "\n":
		return "//"
	default:
		return ""
	}
}

type rawMapping struct {
	Letters          map[string]string `json:"letters"`
	Accented_letters map[string]string `json:"accented_letters"`
	Digits           map[string]string `json:"digits"`
	Punctuations     map[string]string `json:"punctuations"`
}

func NewMapper(jsonFilename string) (Mapper, error) {
	return load(jsonFilename)
}

func load(jsonFilename string) (Mapper, error) {
	fileContent, err := os.ReadFile(jsonFilename)
	if err != nil {
		return nil, err
	}

	var mapping rawMapping

	if err := json.Unmarshal(fileContent, &mapping); err != nil {
		return nil, err
	}

	return &Mapping{mapping: mapping}, nil
}

func createSymbolToMorseTable(m rawMapping) map[string]string {

	table := m.Letters
	maps.Copy(table, m.Accented_letters)
	maps.Copy(table, m.Digits)
	maps.Copy(table, m.Punctuations)

	return table
}

func createMorseToSymbolTable(m rawMapping) map[string]string {

	reverseMapKeyValue := func(m map[string]string) map[string]string {

		output := map[string]string{}

		for k, v := range m {
			output[v] = k
		}

		return output
	}

	table := reverseMapKeyValue(m.Letters)
	maps.Copy(table, reverseMapKeyValue(m.Accented_letters))
	maps.Copy(table, reverseMapKeyValue(m.Digits))
	maps.Copy(table, reverseMapKeyValue(m.Punctuations))

	return table
}
