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
	file                *os.File
	reader              *bufio.Reader
	buffer              []byte
	incompleteUtf8Chunk []byte
}

func (fr *FileReader) trimOffInvalidUtf8Tail(chunk []byte) []byte {
	for !utf8.Valid(chunk) && len(chunk) > 0 {
		invalidBytePosition := len(chunk)
		fr.incompleteUtf8Chunk = append([]byte{chunk[invalidBytePosition-1]}, fr.incompleteUtf8Chunk...)
		chunk = chunk[:invalidBytePosition-1]
	}

	return chunk
}

func (fr *FileReader) ReadChunk() (string, error) {

	n, err := fr.reader.Read(fr.buffer)
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			if len(fr.incompleteUtf8Chunk) > 0 {
				// Discarding incomplete UTF-8 sequence at the end-of-file
				fr.incompleteUtf8Chunk = nil
			}
		}

		return "", err
	}

	chunk := append(fr.incompleteUtf8Chunk, fr.buffer[:n]...)

	fr.incompleteUtf8Chunk = []byte{}

	chunk = fr.trimOffInvalidUtf8Tail(chunk)

	return string(chunk), nil
}

func (fr *FileReader) Close() error {
	if fr.file != nil {
		return fr.file.Close()
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
