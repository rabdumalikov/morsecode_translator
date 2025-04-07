package main

import (
	"flag"
	"fmt"
	"io"
	"morse_converter/morse/rw"
	"os"
	"regexp"
	"strings"
)

type ComparisonResult struct {
	Differences []Diff
}

type Diff struct {
	Line      int
	LineText1 string
	LineText2 string
}

func CompareFilesInChunks(reader1 rw.Reader, reader2 rw.Reader) (ComparisonResult, error) {

	result := ComparisonResult{}
	lineNumber := 1

	file1BufferRemainder := ""
	file2BufferRemainder := ""

	var whitespaceRegex = regexp.MustCompile(`\s+`)

	for {
		chunk1Buf, err1 := reader1.ReadChunk()
		chunk2Buf, err2 := reader2.ReadChunk()

		if err1 != nil && err1 != io.EOF && err1 != io.ErrUnexpectedEOF {
			return result, fmt.Errorf("error reading file1: %v", err1)
		}
		if err2 != nil && err2 != io.EOF && err2 != io.ErrUnexpectedEOF {
			return result, fmt.Errorf("error reading file2: %v", err2)
		}

		// Trim and convert to lower case
		chunk1Normalized := file1BufferRemainder + whitespaceRegex.ReplaceAllString(strings.ToLower(chunk1Buf), "")
		chunk2Normalized := file2BufferRemainder + whitespaceRegex.ReplaceAllString(strings.ToLower(chunk2Buf), "")

		// Reset
		file1BufferRemainder = ""
		file2BufferRemainder = ""

		minLength := min(len(chunk1Normalized), len(chunk2Normalized))

		chunk1ToCompare := chunk1Normalized[:minLength]
		chunk2ToCompare := chunk2Normalized[:minLength]

		if len(chunk1Normalized) > minLength {
			file1BufferRemainder = chunk1Normalized[minLength:]
		}
		if len(chunk2Normalized) > minLength {
			file2BufferRemainder = chunk2Normalized[minLength:]
		}

		if chunk1ToCompare != chunk2ToCompare {
			result.Differences = append(result.Differences, Diff{LineText1: chunk1ToCompare, LineText2: chunk2ToCompare, Line: lineNumber})
		}

		if err1 == io.EOF || err2 == io.EOF {
			break
		}

		lineNumber++
	}

	return result, nil
}

func printResult(results ComparisonResult) {
	for i, diff := range results.Differences {
		fmt.Printf("Difference #%d at line=%d: \nfile1_line=[%s]\nfile2_line=[%s]\n", i, diff.Line, diff.LineText1, diff.LineText2)
	}
}

func main() {

	file1path := flag.String("file1", "", "Path to the first file (required)")
	file2path := flag.String("file2", "", "Path to the second file (required)")

	flag.Parse()

	if *file1path == "" || *file2path == "" {
		fmt.Println("Error: Both file paths are required")
		flag.Usage()
		os.Exit(1)
	}

	reader1, err := rw.NewFileReader(*file1path)

	if err != nil {
		fmt.Printf("Error: Cannot open file1: %v\n", err)
		os.Exit(1)
	}
	defer reader1.Close()

	reader2, err := rw.NewFileReader(*file2path)
	if err != nil {
		fmt.Printf("Error: Cannot open file2: %v\n", err)
		os.Exit(1)
	}
	defer reader2.Close()

	result, err := CompareFilesInChunks(reader1, reader2)
	if err != nil {
		fmt.Printf("Error: Comparing files: %v\n", err)
		os.Exit(1)
	}

	printResult(result)
}
