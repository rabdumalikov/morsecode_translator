package mapping

import (
	"encoding/json"
	"maps"
	"os"
)

type Mapper interface {
	SymbolToMorse(symbol string) (string, bool)
	MorseToSymbol(morse string) (string, bool)
}

type Mapping struct {
	mapping              rawMapping
	symbolToMorseMapping map[string]string
	morseToSymbolMapping map[string]string
}

func (p *Mapping) SymbolToMorse(symbol string) (string, bool) {
	if len(p.symbolToMorseMapping) == 0 {
		p.symbolToMorseMapping = createSymbolToMorseTable(p.mapping)
	}

	value, ok := p.symbolToMorseMapping[symbol]
	return value, ok
}

func (p *Mapping) MorseToSymbol(morse string) (string, bool) {
	if len(p.morseToSymbolMapping) == 0 {
		p.morseToSymbolMapping = createMorseToSymbolTable(p.mapping)
	}

	value, ok := p.morseToSymbolMapping[morse]
	return value, ok
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

	reverseMap := func(m map[string]string) map[string]string {

		output := map[string]string{}

		for k, v := range m {
			output[v] = k
		}

		return output
	}

	table := reverseMap(m.Letters)
	maps.Copy(table, reverseMap(m.Accented_letters))
	maps.Copy(table, reverseMap(m.Digits))
	maps.Copy(table, reverseMap(m.Punctuations))

	return table
}
