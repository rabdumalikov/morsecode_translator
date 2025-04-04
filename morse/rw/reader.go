package rw

import (
	"bufio"
	"io"
	"os"
	"unicode/utf8"
)

const (
	ChunkSize = 64000 // 64KB
)

type Reader interface {
	ReadChunk() (string, error)
	Close() error
}

type FileReader struct {
	file            *os.File
	reader          *bufio.Reader
	buffer          []byte
	bufferRemainder []byte
}

func (p *FileReader) ReadChunk() (string, error) {

	n, err := p.reader.Read(p.buffer)
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			if len(p.bufferRemainder) > 0 {
				chunk := p.bufferRemainder
				p.bufferRemainder = nil
				return string(chunk), err
			}
		}

		return "", err
	}

	chunk := append(p.bufferRemainder, p.buffer[:n]...)

	p.bufferRemainder = []byte{}

	// Trim off invalid UTF-8 tail
	for !utf8.Valid(chunk) && len(chunk) > 0 {
		p.bufferRemainder = append([]byte{chunk[len(chunk)-1]}, p.bufferRemainder...)
		chunk = chunk[:len(chunk)-1]
	}

	return string(chunk), nil
}

func (p *FileReader) Close() error {
	if p.file != nil {
		return p.file.Close()
	}

	return nil
}

func NewFileReader(filename string) (Reader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)

	return &FileReader{file: file, reader: reader, buffer: make([]byte, ChunkSize)}, nil
}
