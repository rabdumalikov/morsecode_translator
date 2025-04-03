package morse

import (
	"errors"
	"fmt"
	"morse_converter/morse/encoding"
	"morse_converter/morse/io"
	"morse_converter/morse/mapping"
	"path/filepath"
)

type TransformationType int

const (
	DecodeMorseCode TransformationType = iota
	EncodeMorseCode
	UnknownTransformation
)

type Converter struct {
	reader             io.Reader
	writer             io.Writer
	transformationType TransformationType
	encoder            encoding.Encoder
	decoder            encoding.Decoder
}

func New(inputFilename, outputFilename string) (*Converter, error) {
	transformationType, err := extensionToTransformationType(filepath.Ext(inputFilename))

	if err != nil {
		return nil, err
	}

	mappings, err := mapping.NewMapper("morse/mapping/symbol2morse.json")
	if err != nil {
		return nil, fmt.Errorf("loading morse mappings: %w", err)
	}

	encoder := encoding.NewEncoder(mappings)
	decoder := encoding.NewDecoder(mappings)

	reader, err := io.NewFileReader(inputFilename, transformationType == DecodeMorseCode)
	if err != nil {
		return nil, fmt.Errorf("creating reader: %w", err)
	}

	writer, err := io.NewFileWriter(outputFilename)
	if err != nil {
		return nil, fmt.Errorf("creating writer: %w", err)
	}

	return &Converter{
		reader:             reader,
		writer:             writer,
		transformationType: transformationType,
		encoder:            *encoder,
		decoder:            *decoder,
	}, nil
}

func (p *Converter) Close() error {
	var errs []error

	if err := p.reader.Close(); err != nil {
		errs = append(errs, err)
	}

	if err := p.writer.Close(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing resources: %v", errs)
	}
	return nil
}

func (p *Converter) Process() error {
	switch p.transformationType {
	case DecodeMorseCode:
		return p.decodeMorseToText()
	case EncodeMorseCode:
		return p.encodeTextToMorse()
	default:
		return errors.New("unknown transformation type")
	}
}

func (p *Converter) decodeMorseToText() error {
	for {
		line, err := p.reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		text := p.decoder.Decode(line) + "\n"
		if err := p.writer.WriteLine(text); err != nil {
			return err
		}
	}

	return nil
}

func (p *Converter) encodeTextToMorse() error {
	for {
		line, err := p.reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		morseCode := p.encoder.Encode(line) + "//"
		if err := p.writer.WriteLine(morseCode); err != nil {
			return err
		}
	}

	return nil
}

func extensionToTransformationType(ext string) (TransformationType, error) {
	switch ext {
	case ".morse":
		return DecodeMorseCode, nil
	case ".txt":
		return EncodeMorseCode, nil
	default:
		return UnknownTransformation, fmt.Errorf("unknown file extension for transformation: %s", ext)
	}
}
