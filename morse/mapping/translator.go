package mapping

import (
	"encoding/json"
	"maps"
	"os"
)

type Translator interface {
	CharToMorse(char string) (string, bool)
	MorseToChar(morse string) (string, bool)
}

type rawMapping struct {
	mapping            jsonMapping
	charToMorseMapping map[string]string
	morseToCharMapping map[string]string
}

func (rm *rawMapping) CharToMorse(char string) (string, bool) {
	if rm.charToMorseMapping == nil {
		rm.charToMorseMapping = createSymbolToMorseTable(rm.mapping)
	}

	value, ok := rm.charToMorseMapping[char]
	return value, ok
}

func (rm *rawMapping) MorseToChar(morse string) (string, bool) {
	if rm.morseToCharMapping == nil {
		rm.morseToCharMapping = createMorseToSymbolTable(rm.mapping)
	}

	value, ok := rm.morseToCharMapping[morse]
	return value, ok
}

type jsonMapping struct {
	Letters          map[string]string `json:"letters"`
	Accented_letters map[string]string `json:"accented_letters"`
	Digits           map[string]string `json:"digits"`
	Punctuations     map[string]string `json:"punctuations"`
}

func NewTranslator(jsonFilename string) (Translator, error) {
	return load(jsonFilename)
}

func load(jsonFilename string) (Translator, error) {
	fileContent, err := os.ReadFile(jsonFilename)
	if err != nil {
		return nil, err
	}

	var mapping jsonMapping

	if err := json.Unmarshal(fileContent, &mapping); err != nil {
		return nil, err
	}

	return &rawMapping{mapping: mapping}, nil
}

func createSymbolToMorseTable(m jsonMapping) map[string]string {

	table := m.Letters
	maps.Copy(table, m.Accented_letters)
	maps.Copy(table, m.Digits)
	maps.Copy(table, m.Punctuations)

	return table
}

func createMorseToSymbolTable(m jsonMapping) map[string]string {

	reverseMapKeyValues := func(m map[string]string) map[string]string {

		output := map[string]string{}

		for k, v := range m {
			output[v] = k
		}

		return output
	}

	table := reverseMapKeyValues(m.Letters)
	maps.Copy(table, reverseMapKeyValues(m.Accented_letters))
	maps.Copy(table, reverseMapKeyValues(m.Digits))
	maps.Copy(table, reverseMapKeyValues(m.Punctuations))

	return table
}
