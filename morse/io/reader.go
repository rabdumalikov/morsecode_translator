package io

import (
	"bufio"
	"bytes"
	"errors"
	"os"
)

var (
	EOF = errors.New("end of file")
)

type Reader interface {
	ReadLine() (string, error)
	Close() error
}

type FileReader struct {
	file    *os.File
	scanner *bufio.Scanner
}

func (p *FileReader) ReadLine() (string, error) {
	if !p.scanner.Scan() {
		if err := p.scanner.Err(); err != nil {
			return "", err
		}
		return "", EOF
	}

	return p.scanner.Text(), nil
}

func (p *FileReader) Close() error {
	if p.file != nil {
		return p.file.Close()
	}

	return nil
}

func NewFileReader(filename string, isMorse bool) (Reader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	if isMorse {
		scanner.Split(scanMorseEndLine)
	} else {
		scanner.Split(scanTxtEndLine)
	}

	return &FileReader{file: file, scanner: scanner}, nil
}

func scanDelimitedLine(data []byte, atEOF bool, separator []byte) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, separator); i >= 0 {
		// We have a full separator-terminated line.
		return i + len(separator), data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

// Custom because we want to preserve \r, so that restored file was identical to original
func scanTxtEndLine(data []byte, atEOF bool) (advance int, token []byte, err error) {
	return scanDelimitedLine(data, atEOF, []byte("\n"))
}

func scanMorseEndLine(data []byte, atEOF bool) (advance int, token []byte, err error) {
	return scanDelimitedLine(data, atEOF, []byte("//"))
}
