package tests

import (
	"morse_converter/morse/mapping"
	"testing"
)

func TestNewLineSeparator(t *testing.T) {

	separator := mapping.NewMorseSeparator(mapping.NewLineSeparator)

	stringRepresentation := separator.ToString()
	if stringRepresentation != "//" {
		t.Errorf("Wrong morse new line representation got [%s]; want [//]", stringRepresentation)
	}

	stringRepresentation = separator.ToText().ToString()
	if stringRepresentation != "\n" {
		t.Errorf("Wrong text new line representation got [%s]; want [\\n]", stringRepresentation)
	}
}

func TestWordSeparator(t *testing.T) {

	separator := mapping.NewMorseSeparator(mapping.WordSeparator)

	stringRepresentation := separator.ToString()
	if stringRepresentation != "/" {
		t.Errorf("Wrong morse new line representation got [%s]; want [/]", stringRepresentation)
	}

	stringRepresentation = separator.ToText().ToString()
	if stringRepresentation != " " {
		t.Errorf("Wrong text new line representation got [%s]; want [ ]", stringRepresentation)
	}
}

func TestCharSeparator(t *testing.T) {

	separator := mapping.NewMorseSeparator(mapping.CharSeparator)

	stringRepresentation := separator.ToString()
	if stringRepresentation != " " {
		t.Errorf("Wrong morse new line representation got [%s]; want [ ]", stringRepresentation)
	}

	stringRepresentation = separator.ToText().ToString()
	if stringRepresentation != "" {
		t.Errorf("Wrong text new line representation got [%s]; want []", stringRepresentation)
	}
}

func TestNewLineSeparatorReversed(t *testing.T) {

	separator := mapping.NewTextSeparator(mapping.NewLineSeparator)

	stringRepresentation := separator.ToString()
	if stringRepresentation != "\n" {
		t.Errorf("Wrong text new line representation got [%s]; want [\\n]", stringRepresentation)
	}

	stringRepresentation = separator.ToMorse().ToString()
	if stringRepresentation != "//" {
		t.Errorf("Wrong morse new line representation got [%s]; want [//]", stringRepresentation)
	}
}

func TestWordSeparatorReversed(t *testing.T) {

	separator := mapping.NewTextSeparator(mapping.WordSeparator)

	stringRepresentation := separator.ToString()
	if stringRepresentation != " " {
		t.Errorf("Wrong text new line representation got [%s]; want [ ]", stringRepresentation)
	}

	stringRepresentation = separator.ToMorse().ToString()
	if stringRepresentation != "/" {
		t.Errorf("Wrong morse new line representation got [%s]; want [/]", stringRepresentation)
	}
}

func TestCharSeparatorReversed(t *testing.T) {

	separator := mapping.NewTextSeparator(mapping.CharSeparator)

	stringRepresentation := separator.ToString()
	if stringRepresentation != "" {
		t.Errorf("Wrong text new line representation got [%s]; want []", stringRepresentation)
	}

	stringRepresentation = separator.ToMorse().ToString()
	if stringRepresentation != " " {
		t.Errorf("Wrong morse new line representation got [%s]; want [ ]", stringRepresentation)
	}
}
func TestSeparatorConstants(t *testing.T) {

	if mapping.MorseWordSeparator != '/' {
		t.Errorf("Wrong MorseWordSeparator representation got [%c]; want [/]", mapping.MorseWordSeparator)
	}

	if mapping.MorseCharSeparator != ' ' {
		t.Errorf("Wrong MorseCharSeparator representation got [%c]; want [ ]", mapping.MorseCharSeparator)
	}

	if mapping.MorseNewLine != "//" {
		t.Errorf("Wrong MorseNewLine representation got [%s]; want [//]", mapping.MorseNewLine)
	}

	if mapping.TextWordSeparator != ' ' {
		t.Errorf("Wrong TextWordSeparator representation got [%c]; want [ ]", mapping.TextWordSeparator)
	}

	if mapping.TextCharSeparator != "" {
		t.Errorf("Wrong TextCharSeparator representation got [%s]; want []", mapping.TextCharSeparator)
	}

	if mapping.TextNewLine != '\n' {
		t.Errorf("Wrong TextNewLine representation got [%c]; want [\\n]", mapping.TextNewLine)
	}
}
