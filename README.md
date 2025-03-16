# renafont

CLI tool to rename `Font Family` in TTF files.

## Description

`renafont` is a simple command-line tool written in Go that allows you to rename the `Font Family` name in TTF (TrueType Font) files.

## Features

- Rename `Font Family` names in TTF files
- Easy to use CLI interface

## Installation

To install the `renafont` tool, you need to have Go installed on your system. Then you can use the following command to install it:

```sh
go install github.com/asalvi0/renafont@latest
```

## Usage

To use the `renafont` tool, simply run the following command:

```sh
renafont -input <input-ttf-file> -output <output-ttf-file> -name <new-font-family-name>
```

### Options

- `-input` : The path to the input TTF file.
- `-output` : The path to the output TTF file.
- `-name` : The new `Font Family` name.

### Example

```sh
renafont -input oldfont.ttf -output newfont.ttf -name "New Font Family"
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the GPL-3.0 license. See the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [Go](https://golang.org/)
- [TrueType Font](https://en.wikipedia.org/wiki/TrueType)
