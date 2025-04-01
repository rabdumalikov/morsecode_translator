package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strings"
)

func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func fileExtCorrect(filename string) bool {
	ext := filepath.Ext(filename)

	switch ext {
	case ".morse", ".txt":
		return true
	default:
		return false
	}
}

const (
	DecodeMorseCode = iota
	EncodeMorseCode
	UnknownTransformation
)

type TransformationType int

// name is not good, change it
type InputHandler struct {
	readFromFD         *os.File
	writeToFD          *os.File
	transformationType TransformationType
	charToMorseMapping map[string]string
	morseToCharMapping map[string]string
}

func extensionToTransformationType(ext string) (TransformationType, error) {
	switch ext {
	case ".morse":
		return DecodeMorseCode, nil
	case ".txt":
		return EncodeMorseCode, nil
	default:
		return UnknownTransformation, errors.New("Unknown Transformation")
	}
}

func createOutputHandler(filename string) *os.File {
	if filename == "" {
		return os.Stdout
	}

	fd, err := os.Create(filename)
	if err != nil {
		fmt.Println("Creating/Truncating output file failed with:", err, "Thus print in console")
		return os.Stdout
	}
	return fd
}

func createInputHandler(inputFilename, outputFilename string) *InputHandler {

	fd, err := os.Open(inputFilename)
	if err != nil {
		fmt.Println("Openning input file failed with:", err)
		return nil
	}

	ext := filepath.Ext(inputFilename)
	transformType, err := extensionToTransformationType(ext)

	if err == nil {
		mapping, _ := parseMorseEncoding()

		charToMorseMapping := createSymbolToMorseTable(mapping)
		morseToCharMapping := createMorseToSymbolTable(mapping)

		return &InputHandler{readFromFD: fd, writeToFD: createOutputHandler(outputFilename),
			transformationType: transformType, charToMorseMapping: charToMorseMapping,
			morseToCharMapping: morseToCharMapping}
	} else {
		fmt.Println("Error:", err)
		return nil
	}
}

func (p *InputHandler) toMorseCode(text string) string {
	words := strings.Split(text, " ")
	var encodedWords []string
	for _, word := range words {
		var encodedWord []string

		for _, char := range word {
			upperChar := strings.ToUpper(string(char))
			value, ok := p.charToMorseMapping[upperChar]
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

func (p *InputHandler) toText(morseCode string) string {
	words := strings.Split(morseCode, "/")
	var decodedWords []string
	for _, word := range words {
		decodedWord := ""

		for _, code := range strings.Split(word, " ") {
			value, ok := p.morseToCharMapping[code]
			if ok {
				decodedWord += value
			} else {
				// TODO: think about how am I handling if symbol not part of table
				decodedWord += code
			}
		}
		decodedWords = append(decodedWords, decodedWord)
	}

	return strings.Join(decodedWords, " ") + "\n"
}

func ScanMorseEndLine(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte("//")); i >= 0 {
		// We have a full newline-terminated line.
		return i + 2, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

func (p *InputHandler) processTransformation() {
	if p.transformationType == DecodeMorseCode {
		fmt.Println("Decoding Morse Code")

		scanner := bufio.NewScanner(p.readFromFD)
		scanner.Split(ScanMorseEndLine)

		for scanner.Scan() {
			line := scanner.Text()
			p.toText(line)
			fmt.Fprintln(p.writeToFD, p.toText(line))
		}
	} else if p.transformationType == EncodeMorseCode {
		fmt.Println("Encoding Morse Code")
		// we have to determine symbols that are not part of our vocabulary
		scanner := bufio.NewScanner(p.readFromFD)
		for scanner.Scan() {
			morseCode := p.toMorseCode(scanner.Text())

			fmt.Fprintln(p.writeToFD, morseCode)
		}
	}
}

type Mapping struct {
	Letters          map[string]string `json:"letters"`
	Accented_letters map[string]string `json:"accented_letters"`
	Digits           map[string]string `json:"digits"`
	Punctuations     map[string]string `json:"punctuations"`
}

// TODO: I think that I should have to define it locally,
// but use file in case if I wanna extend
func parseMorseEncoding() (Mapping, error) {
	jsonFilename := "internal/symbol2morse.json"

	fileContent, err := os.ReadFile(jsonFilename)
	if err != nil {
		return Mapping{}, errors.New("Morse Code library is missing!")
	}

	var m Mapping

	// parse json
	jsonParsingErr := json.Unmarshal(fileContent, &m)

	return m, jsonParsingErr
}

func createSymbolToMorseTable(m Mapping) map[string]string {

	table := m.Letters
	maps.Copy(table, m.Accented_letters)
	maps.Copy(table, m.Digits)
	maps.Copy(table, m.Punctuations)

	return table
}

func createMorseToSymbolTable(m Mapping) map[string]string {

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

func main() {
	inputFilename := flag.String("i", "", "Input file name (.morse or .txt)")
	outputFilename := flag.String("o", "", "Output file name")
	help := flag.Bool("Help", false, "Show help message")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *inputFilename == "" {
		fmt.Println("Input file is required")

		flag.Usage()
		os.Exit(1)
	}

	// validate correct format of input file
	if !fileExtCorrect(*inputFilename) {
		fmt.Printf("Unknown format [%s] of input file, should .morse or .txt\n", *inputFilename)
		os.Exit(1)
	}

	// Validate that the file exist
	if !fileExist(*inputFilename) {
		fmt.Printf("Input file [%s] doesn't exist\n", *inputFilename)
		os.Exit(1)
	}

	inputHandler := createInputHandler(*inputFilename, *outputFilename)

	if inputHandler == nil {
		// decide something upon where to print error, right now it is mixed
		os.Exit(1)
	}

	inputHandler.processTransformation()
}

// TODO: Close files
