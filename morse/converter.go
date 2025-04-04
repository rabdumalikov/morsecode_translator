package morse

import (
	"errors"
	"fmt"
	"io"
	"morse_converter/morse/encoding"
	"morse_converter/morse/mapping"
	"morse_converter/morse/rw"
	"path/filepath"
	"strings"
)

type TransformationType int

const (
	DecodeMorseCode TransformationType = iota
	EncodeMorseCode
	UnknownTransformation
)

type Converter struct {
	reader             rw.Reader
	writer             rw.Writer
	transformationType TransformationType
	encoder            encoding.Encoder
	decoder            encoding.Decoder
	mappings           mapping.Mapper
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

	reader, err := rw.NewFileReader(inputFilename)
	if err != nil {
		return nil, fmt.Errorf("creating reader: %w", err)
	}

	writer, err := rw.NewFileWriter(outputFilename)
	if err != nil {
		return nil, fmt.Errorf("creating writer: %w", err)
	}

	return &Converter{
		reader:             reader,
		writer:             writer,
		transformationType: transformationType,
		encoder:            *encoder,
		decoder:            *decoder,
		mappings:           mappings,
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
	buffer := ""

	for {
		chunk, err := p.reader.ReadChunk()
		if err != nil && err != io.EOF {
			return err
		}

		buffer += chunk

		for {
			morseSeparator := "//"

			sepPos := strings.Index(buffer, morseSeparator) // For line/sentence breaks
			if sepPos == -1 {
				// If no line separator, look for word separator

				morseSeparator = "/"
				sepPos = strings.LastIndex(buffer, morseSeparator)

				if sepPos == -1 {
					break
				}

				// Avoid misinterpreting "//" as a single "/" character
				// This could happen if only the first part of "//" is processed
				if sepPos == len(buffer)-1 {
					break
				}
			}

			morseSegment := buffer[:sepPos]
			buffer = buffer[sepPos+len(morseSeparator):]
			textSeparator := p.mappings.ToTextSeparator(morseSeparator)
			if err := p.writer.WriteChunk(p.decoder.Decode(morseSegment) + textSeparator); err != nil {
				return err
			}
		}

		if err == io.EOF {
			// Process any remaining data
			if buffer != "" {
				if err := p.writer.WriteChunk(p.decoder.Decode(buffer)); err != nil {
					return err
				}
			}
			break
		}
	}

	return nil
}

func (p *Converter) encodeTextToMorse() error {
	buffer := ""

	for {
		chunk, err := p.reader.ReadChunk()
		if err != nil && err != io.EOF {
			return err
		}

		buffer += chunk

		for {
			textSeparator := "\n"
			sepPos := strings.Index(buffer, textSeparator) // For line/sentence breaks
			if sepPos == -1 {
				// If no line separator, look for word separator
				textSeparator = " "
				sepPos = strings.LastIndex(buffer, textSeparator)

				if sepPos == -1 {
					break
				}
			}

			subbuffer := buffer[:sepPos]
			buffer = buffer[sepPos+len(textSeparator):] // Remove processed chunk plus space

			morseCode := p.encoder.Encode(subbuffer) + p.mappings.ToMorseSeparator(textSeparator)
			if err := p.writer.WriteChunk(morseCode); err != nil {
				return err
			}
		}

		if err == io.EOF {
			// Process any remaining data
			if buffer != "" {
				morseCode := p.encoder.Encode(buffer)
				if err := p.writer.WriteChunk(morseCode); err != nil {
					return err
				}
			}
			break
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
