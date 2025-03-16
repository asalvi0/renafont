package main

import (
    "flag"
    "fmt"
    "log"

    . "renafont/internal"
)

func main() {
    filePath := flag.String("file", "", "Path to a TTF file")
    newFamily := flag.String("name", "", "New font family name")
    flag.Parse()

    if *filePath == "" || *newFamily == "" {
        fmt.Println("Usage: go run main.go -file <path/to/font.ttf> -name \"NewFontFamily\"")
        return
    }

    if err := ProcessFontFile(*filePath, *newFamily); err != nil {
        log.Fatalf("Error: %v", err)
    }
}
