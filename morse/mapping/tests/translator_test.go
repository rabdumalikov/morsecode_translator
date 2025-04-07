package tests

import (
	"morse_converter/morse/mapping"
	"strings"
	"testing"
)

func TestNotExistingJson(t *testing.T) {

	translator, err := mapping.NewTranslator("")

	if translator != nil || err == nil {
		t.Errorf("Got non-nil translator=[%p] or err=[%v]; expected nil translator, and non-nil error", translator, err)
	}

	if !strings.Contains(err.Error(), "no such file or directory") {
		t.Errorf("Got wrong error for non-existing file err=[%v]; want [no such file or directory]", err)
	}

	translator, err = mapping.NewTranslator("data/doesNotExist.json")

	if translator != nil || err == nil {
		t.Errorf("Got non-nil translator=[%p] or err=[%v]; expected nil translator, and non-nil error", translator, err)
	}

	if !strings.Contains(err.Error(), "no such file or directory") {
		t.Errorf("Got wrong error for non-existing file err=[%v]; want [no such file or directory]", err)
	}
}

func TestExistingButEmptyJson(t *testing.T) {

	translator, err := mapping.NewTranslator("data/empty.json")

	if translator != nil || err == nil {
		t.Errorf("Got non-nil translator=[%p] or err=[%v]; expected nil translator, and non-nil error", translator, err)
	}

	if !strings.Contains(err.Error(), "unexpected end of JSON input") {
		t.Errorf("Got wrong error for non-existing file err=[%v]; want [unexpected end of JSON input]", err)
	}
}

func TestExistingButMalformedJson(t *testing.T) {

	translator, err := mapping.NewTranslator("data/malformed.json")

	if translator != nil || err == nil {
		t.Errorf("Got non-nil translator=[%p] or err=[%v]; expected nil translator, and non-nil error", translator, err)
	}

	if !strings.Contains(err.Error(), "unexpected end of JSON input") {
		t.Errorf("Got wrong error for non-existing file err=[%v]; want [unexpected end of JSON input]", err)
	}
}

func TestCorrectJson(t *testing.T) {

	translator, err := mapping.NewTranslator("data/correct.json")

	if translator == nil || err != nil {
		t.Errorf("Got nil translator=[%p] or non-nil err=[%v]; expected non-nil translator, and nil error", translator, err)
	}
}
