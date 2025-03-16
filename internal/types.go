package internal

import (
	"bytes"
	"encoding/binary"
)

// TTF structures.
type OffsetTable struct {
	SfntVersion   uint32
	NumTables     uint16
	SearchRange   uint16
	EntrySelector uint16
	RangeShift    uint16
}

type TableRecord struct {
	Tag      [4]byte
	CheckSum uint32
	Offset   uint32
	Length   uint32
}

type NameTableHeader struct {
	Format       uint16
	Count        uint16
	StringOffset uint16
}

type NameRecord struct {
	PlatformID uint16
	EncodingID uint16
	LanguageID uint16
	NameID     uint16
	Length     uint16
	Offset     uint16
}

type Table struct {
	Tag  string
	Data []byte
}

// readStruct is a generic helper to read a structure from a bytes.Reader.
func readStruct[T any](r *bytes.Reader) (T, error) {
	var t T
	return t, binary.Read(r, binary.BigEndian, &t)
}

// writeStructs writes a slice of values to the given buffer using binary.BigEndian.
func writeStructs[T any](buf *bytes.Buffer, slice []T) error {
	for _, v := range slice {
		if err := binary.Write(buf, binary.BigEndian, v); err != nil {
			return err
		}
	}
	return nil
}
