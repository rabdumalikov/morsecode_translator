package rw

import (
	"fmt"
	"os"
)

type Writer interface {
	WriteChunk(chunk string) error
	Close() error
}

type FileWriter struct {
	file *os.File
}

func (p *FileWriter) WriteChunk(chunk string) error {

	_, err := fmt.Fprint(p.file, chunk)

	return err
}

func (p *FileWriter) Close() error {
	if p.file != nil && p.file != os.Stdout {
		return p.file.Close()
	}

	return nil
}

func NewFileWriter(filename string) (Writer, error) {

	if filename == "" {
		return &FileWriter{file: os.Stdout}, nil
	}

	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return &FileWriter{file: file}, nil
}
