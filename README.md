# Morse Code Translator

## Overview
The Morse Code Translator is a powerful command-line tool that converts between text and Morse code formats. Designed with performance in mind, it efficiently processes files of any size through optimized chunk-based reading, ensuring smooth operation even with extensive data.

## Features

- **Text to Morse Conversion**: Transform plain text files into Morse code with precision
  - Characters without Morse code representations are preserved as-is in the output

- **Morse to Text Conversion**: Decode Morse code back into readable text
  - Characters without text equivalents are maintained unchanged in the resulting file

- **Advanced Chunk Processing**: Handles large files efficiently with 64KB chunk processing
  - Enhanced processing respects content structure by splitting chunks at natural newline boundaries

- **Flexible Output Options**: Direct results to custom file paths or standard output

- **Expandable Character Support**: Customize the Morse code dictionary by adding or removing character representations

## Morse Code Mappings



## Installation

Ensure you have Go >=1.24.1 installed on your system. Clone the repository and navigate to the project directory:

```bash
git clone https://github.com/rabdumalikov/morsecode_translator.git
cd morsecode_translator/
```

Build the project using:

```bash
go build -o translator
```

## Usage

Run the morsecode translator with the following command:

```bash
./translator -i <input-file> -o <output-file>
```

or 

```bash
./translator -i <input-file>
```


- `-i`: Specifies the input file path (required).
- `-o`: Specifies the output file path (optional, defaults to stdout).
- `-h`: Displays the help message.

## Compare Files:

Run the comparator tool with the following command:

```bash
go run tools/comparator.go -file1 <first-file> -file2 <second-file>
```


Alternatively, use linux native tool **diff**:

```bash
diff -iw <first-file> <second-file>
```

- `-i`: case insensitive comparison.
- `-w`: ignoring all white-spaces.

## Morse Code Extension:

To support additional characters in Morse code, modify the following [JSON](morse/mapping/char2morse.json) accordingly.

## Performance

The morsecode translator performance metrics:

- **Text to Morse**: Processes 500MB text with newlines in approximately 22.598 seconds.
- **Morse to Text**: Decodes 500MB Morse code in approximately 32.895 seconds.

## Project Structure

- `main.go`: The entry point of the application, handling command-line arguments and initiating the translation process.
- `morse/`: Contains the core logic for encoding and decoding Morse code.
  - `converter.go`: Implements the conversion logic.
  - `encoding/`: Handles the encoding and decoding processes.
  - `mapping/`: Manages character to Morse code mappings.
  - `rw/`: Provides utilities for reading and writing files.
- `tools/`: Provides application to compare two text files.

## License

This project is licensed under the MIT License.
