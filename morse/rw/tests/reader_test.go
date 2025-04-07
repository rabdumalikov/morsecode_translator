package tests

import (
	"io"
	"morse_converter/morse/rw"
	"os"
	"strings"
	"testing"
)

func createTempFile(content string) (*os.File, error) {
	tmpFile, err := os.CreateTemp("", "test_file")
	if err != nil {
		return nil, err
	}
	_, err = tmpFile.WriteString(content)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return nil, err
	}
	_, err = tmpFile.Seek(0, io.SeekStart) // Reset file pointer to the beginning
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return nil, err
	}
	return tmpFile, nil
}

func TestSingleReadChunk(t *testing.T) {
	expectedContent := "This is a test file content."
	tmpFile, err := createTempFile(expectedContent)
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	reader, err := rw.NewFileReader(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewFileReader failed: %v", err)
		return
	}
	defer reader.Close()

	chunk, err := reader.ReadChunk()
	if err != nil {
		t.Fatalf("ReadChunk failed: %v", err)
		return
	}

	if chunk != expectedContent {
		t.Errorf("ReadChunk returned incorrect content: got '%s', want '%s'", chunk, expectedContent)
		return
	}

	_, err = reader.ReadChunk()
	if err != io.EOF {
		t.Errorf("Expected io.EOF after reading all content, got: %v", err)
	}
}

func TestMultipleReadChunks(t *testing.T) {
	expectedContent := strings.Repeat("A", rw.ChunkSize)
	expectedContent += strings.Repeat("B", rw.ChunkSize)
	expectedContent += strings.Repeat("C", 100)

	tmpFile, err := createTempFile(expectedContent)
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	reader, err := rw.NewFileReader(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewFileReader failed: %v", err)
		return
	}
	defer reader.Close()

	var readContent string
	for {
		chunk, err := reader.ReadChunk()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("ReadChunk failed: %v", err)
			return
		}
		readContent += chunk
	}

	if readContent != expectedContent {
		t.Errorf("ReadChunk did not read the entire content correctly. Read %d bytes, expected %d bytes.", len(readContent), len(expectedContent))
	}
}

func TestIncompleteUTF8(t *testing.T) {
	// '€' is 3 bytes in UTF-8 (E2 82 AC)
	incompleteEuroSign := "\xE2"

	expectedChunk1 := strings.Repeat("x", rw.ChunkSize-len(incompleteEuroSign))
	expectedChunk2 := "\xE2\x82\xAC"

	tmpFile, err := createTempFile(expectedChunk1 + expectedChunk2)
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	reader, err := rw.NewFileReader(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewFileReader failed: %v", err)
		return
	}
	defer reader.Close()

	chunk1, err := reader.ReadChunk()
	if err != nil {
		t.Fatalf("ReadChunk (first) failed: %v", err)
		return
	}

	chunk2, err := reader.ReadChunk()
	if err != nil {
		t.Fatalf("ReadChunk (second) failed: %v", err)
		return
	}

	if chunk1 != expectedChunk1 {
		t.Errorf("ReadChunk (first) returned incorrect content: got '%s', want '%s'", chunk1, expectedChunk1)
		return
	}
	if chunk2 != expectedChunk2 {
		t.Errorf("ReadChunk (second) returned incorrect content: got '%s', want '%s'", chunk2, expectedChunk2)
	}
}

func TestIncompleteUTF8Extended(t *testing.T) {
	// '€' is 3 bytes in UTF-8 (E2 82 AC)
	incompleteEuroSign := "\xE2\x82"

	expectedChunk1 := strings.Repeat("x", rw.ChunkSize-len(incompleteEuroSign))
	expectedChunk2 := "\xE2\x82\xAC"

	tmpFile, err := createTempFile(expectedChunk1 + expectedChunk2)
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	reader, err := rw.NewFileReader(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewFileReader failed: %v", err)
		return
	}
	defer reader.Close()

	chunk1, err := reader.ReadChunk()
	if err != nil {
		t.Fatalf("ReadChunk (first) failed: %v", err)
		return
	}

	chunk2, err := reader.ReadChunk()
	if err != nil {
		t.Fatalf("ReadChunk (second) failed: %v", err)
		return
	}

	if chunk1 != expectedChunk1 {
		t.Errorf("ReadChunk (first) returned incorrect content: got '%s', want '%s'", chunk1, expectedChunk1)
		return
	}
	if chunk2 != expectedChunk2 {
		t.Errorf("ReadChunk (second) returned incorrect content: got '%s', want '%s'", chunk2, expectedChunk2)
	}
}

func TestDiscardingIncompleteUTF8(t *testing.T) {
	// '€' is 3 bytes in UTF-8 (E2 82 AC)
	incompleteEuroSign := "\xE2\x82"
	expectedChunk1 := strings.Repeat("x", rw.ChunkSize-len(incompleteEuroSign))
	expectedChunk2 := ""

	tmpFile, err := createTempFile(expectedChunk1 + incompleteEuroSign)
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	reader, err := rw.NewFileReader(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewFileReader failed: %v", err)
		return
	}
	defer reader.Close()

	chunk, err := reader.ReadChunk()
	if err != nil {
		t.Fatalf("ReadChunk failed: %v", err)
		return
	}

	if chunk != expectedChunk1 {
		t.Errorf("ReadChunk returned incorrect content (incomplete UTF-8 at EOF): got '%s', want '%s'", chunk, expectedChunk1)
		return
	}

	chunk, err = reader.ReadChunk()
	if err != io.EOF {
		t.Errorf("Expected io.EOF after reading all content (with incomplete UTF-8), got: %v", err)
		return
	}

	if chunk != expectedChunk2 {
		t.Errorf("ReadChunk returned incorrect content (incomplete UTF-8 at EOF): got '%s', want '%s'", chunk, expectedChunk2)
	}
}
