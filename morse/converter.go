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
}

func New(inputFilename, outputFilename string) (*Converter, error) {
	transformationType, err := extensionToTransformationType(filepath.Ext(inputFilename))

	if err != nil {
		return nil, err
	}

	mappings, err := mapping.NewTranslator("morse/mapping/char2morse.json")
	if err != nil {
		return nil, fmt.Errorf("loading morse mappings: %w", err)
	}

	encoder := encoding.NewEncoder(mappings)
	decoder := encoding.NewDecoder(mappings)

	reader, err := rw.NewFileReader(inputFilename)
	if err != nil {
		return nil, fmt.Errorf("creating reader for [%s]: %w", inputFilename, err)
	}

	writer, err := rw.NewFileWriter(outputFilename)
	if err != nil {
		return nil, fmt.Errorf("creating writer for [%s]: %w", outputFilename, err)
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
	buffer := ""

	for {
		chunk, err := p.reader.ReadChunk()
		if err != nil && err != io.EOF {
			return err
		}

		buffer += chunk

		for {
			separator := mapping.NewMorseSeparator(mapping.NewLineSeparator)

			sepPos := strings.Index(buffer, separator.ToString()) // For line/sentence breaks
			if sepPos == -1 {
				// If no line separator, look for word separator
				separator = mapping.NewMorseSeparator(mapping.WordSeparator)
				sepPos = strings.LastIndex(buffer, separator.ToString())

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
			buffer = buffer[sepPos+len(separator.ToString()):]

			textSeparator := separator.ToText().ToString()

			textOutput, decodeLineErr := p.decoder.DecodeLine(morseSegment)

			if decodeLineErr != nil {
				return decodeLineErr
			}
			if err := p.writer.WriteChunk(textOutput + textSeparator); err != nil {
				return err
			}
		}

		if err == io.EOF {
			// Process any remaining data
			if buffer != "" {
				textOutput, decodeLineErr := p.decoder.DecodeLine(buffer)

				if decodeLineErr != nil {
					return decodeLineErr
				}

				if err := p.writer.WriteChunk(textOutput); err != nil {
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
			separator := mapping.NewTextSeparator(mapping.NewLineSeparator)

			sepPos := strings.Index(buffer, separator.ToString()) // For line/sentence breaks
			if sepPos == -1 {
				// If no line separator, look for word separator
				separator = mapping.NewTextSeparator(mapping.WordSeparator)
				sepPos = strings.LastIndex(buffer, separator.ToString())

				if sepPos == -1 {
					break
				}
			}

			textSegment := buffer[:sepPos]
			buffer = buffer[sepPos+len(separator.ToString()):] // Remove processed chunk plus space
			morseSeparator := separator.ToMorse().ToString()

			morseOutput, encodeLineErr := p.encoder.EncodeLine(textSegment)
			if encodeLineErr != nil {
				return encodeLineErr
			}

			if err := p.writer.WriteChunk(morseOutput + morseSeparator); err != nil {
				return err
			}
		}

		if err == io.EOF {
			// Process any remaining data
			if buffer != "" {
				morseOutput, encodeLineErr := p.encoder.EncodeLine(buffer)

				if encodeLineErr != nil {
					return encodeLineErr
				}

				if err := p.writer.WriteChunk(morseOutput); err != nil {
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
