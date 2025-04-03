package main

import (
	"bufio"
	"flag"
	"fmt"
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

func CompareFiles(file1 *os.File, file2 *os.File) (ComparisonResult, error) {
	scanner1 := bufio.NewScanner(file1)
	scanner2 := bufio.NewScanner(file2)

	result := ComparisonResult{}

	lineNumber := 1
	for scanner1.Scan() && scanner2.Scan() {
		var whitespaceRegex = regexp.MustCompile(`\s{2,}`)

		// Remove all whitespace characters from both lines
		line1 := whitespaceRegex.ReplaceAllString(strings.TrimSpace(strings.ToLower(scanner1.Text())), " ")
		line2 := whitespaceRegex.ReplaceAllString(strings.TrimSpace(strings.ToLower(scanner2.Text())), " ")

		if len(line1) != len(line2) && line1 != line2 {
			result.Differences = append(result.Differences, Diff{LineText1: line1, LineText2: line2, Line: lineNumber})
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

	file1Path := flag.String("file1", "", "Path to the first file (required)")
	file2Path := flag.String("file2", "", "Path to the second file (required)")

	flag.Parse()

	if *file1Path == "" || *file2Path == "" {
		fmt.Println("Error: Both file paths are required")
		flag.Usage()
		os.Exit(1)
	}

	file1, err := os.Open(*file1Path)
	if err != nil {
		fmt.Printf("Error: Cannot open file1: %v", err)
		os.Exit(1)
	}

	defer file1.Close()

	file2, err := os.Open(*file2Path)
	if err != nil {
		fmt.Printf("Error: Cannot open file2: %v", err)
		os.Exit(1)
	}

	defer file2.Close()

	result, err := CompareFiles(file1, file2)
	if err != nil {
		fmt.Printf("Error: comparing files: %v", err)
		os.Exit(1)
	}

	printResult(result)
}
