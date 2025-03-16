package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"seehuhn.de/go/sfnt"
)

const (
	supportedExt = ".ttf, .otf, .woff, .woff2"
	usageMessage = "Usage: renafont -in <path/to/input.ttf> -name \"NewFontFamily\"\nSupported formats: " + supportedExt
)

// Returns true if the filename has a supported font extension
func isFontFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return strings.Contains(supportedExt, ext)
}

func main() {
	inFilePath := flag.String("in", "", "Path to a font file")
	newFamily := flag.String("name", "", "New font family name")
	flag.Parse()

	if *inFilePath == "" || *newFamily == "" {
		fmt.Println(usageMessage)
		return
	}

	if !isFontFile(*inFilePath) {
		log.Fatalf("unsupported font format: %s", filepath.Ext(*inFilePath))
	}

	font, err := sfnt.ReadFile(*inFilePath)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	font.FamilyName = *newFamily

	// Get file extension
	ext := strings.ToLower(filepath.Ext(*inFilePath))

	// Get file name without extension
	fileName := strings.TrimSuffix(*inFilePath, ext)
	outFilePath := fileName + ".modified" + ext

	// Create an output file for the modified font
	out, err := os.Create(outFilePath)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer out.Close()

	// Write the modified font data
	n, err := font.Write(out)
	if err != nil {
		log.Fatalf("Error writing font: %v", err)
	}
	log.Printf("Wrote %d bytes to %s", n, outFilePath)
}
