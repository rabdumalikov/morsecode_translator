package tests

import (
	"bytes"
	"io"
	"morse_converter/morse/rw"
	"os"
	"strings"
	"testing"
)

func createFile(filename string) (*os.File, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func readTestFile() (string, error) {

	data, err := os.ReadFile("test_file")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func TestSingleWriteChunk(t *testing.T) {
	expectedContent := "This is a test file content."
	file, err := createFile("test_file")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
		return
	}
	defer os.Remove(file.Name())

	writer, err := rw.NewFileWriter(file.Name())
	if err != nil {
		t.Fatalf("NewFileWriter failed: %v", err)
		return
	}
	defer writer.Close()

	err = writer.WriteChunk(expectedContent)
	if err != nil {
		t.Fatalf("WriteChunk failed: %v", err)
		return
	}

	actualContent, err := readTestFile()

	if err != nil {
		t.Fatalf("readTestFile failed: %v", err)
		return
	}

	if actualContent != expectedContent {
		t.Errorf("Expected io.EOF after reading all content, got: %v", err)
	}
}

func TestMultipleWriteChunk(t *testing.T) {
	file, err := createFile("test_file")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
		return
	}
	defer os.Remove(file.Name())

	writer, err := rw.NewFileWriter(file.Name())
	if err != nil {
		t.Fatalf("NewFileWriter failed: %v", err)
		return
	}
	defer writer.Close()

	expectedContent := []string{}

	expectedContent = append(expectedContent, strings.Repeat("A", rw.ChunkSize))
	expectedContent = append(expectedContent, strings.Repeat("B", rw.ChunkSize))
	expectedContent = append(expectedContent, strings.Repeat("C", 100))

	for _, content := range expectedContent {
		err = writer.WriteChunk(content)
		if err != nil {
			t.Fatalf("WriteChunk failed: %v", err)
			return
		}
	}

	actualContent, err := readTestFile()
	if err != nil {
		t.Fatalf("readTestFile failed: %v", err)
		return
	}

	if actualContent != strings.Join(expectedContent, "") {
		t.Errorf("WriteChunk did not write the entire content correctly. Read %d bytes, expected %d bytes.", len(actualContent), len(strings.Join(expectedContent, "")))
	}
}

func TestSingleWriteChunkToStdout(t *testing.T) {
	expectedContent := "This is a test file content."

	// Capture stdout
	originalStdout := os.Stdout
	defer func() { os.Stdout = originalStdout }()

	pipeReader, pipeWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
		return
	}
	os.Stdout = pipeWriter

	writer, err := rw.NewFileWriter("")
	if err != nil {
		t.Fatalf("NewFileWriter failed: %v", err)
		return
	}
	defer writer.Close()

	err = writer.WriteChunk(expectedContent)
	if err != nil {
		t.Fatalf("WriteChunk failed: %v", err)
		return
	}

	err = pipeWriter.Close()
	if err != nil {
		t.Fatalf("Failed to close writer: %v", err)
		return
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, pipeReader)
	if err != nil {
		t.Fatalf("Failed to read from reader: %v", err)
		return
	}

	actualStdout := buf.String()

	// Assert against the captured stdout
	if actualStdout != expectedContent {
		t.Errorf("Expected stdout output '%s', got '%s'", expectedContent, actualStdout)
	}
}
