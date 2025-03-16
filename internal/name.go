package internal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"unicode/utf16"
)

func rebuildNameTable(data []byte, newFamily string) ([]byte, error) {
	r := bytes.NewReader(data)
	header, err := readStruct[NameTableHeader](r)
	if err != nil {
		return nil, err
	}

	count := int(header.Count)
	records := make([]NameRecord, count)
	for i := 0; i < count; i++ {
		rec, err := readStruct[NameRecord](r)
		if err != nil {
			return nil, err
		}
		records[i] = rec
	}

	stringStorage := data[header.StringOffset:]
	var newBlock bytes.Buffer
	newRecords := make([]NameRecord, count)
	for i, rec := range records {
		newRec := rec
		origBytes := []byte{}
		if int(rec.Offset)+int(rec.Length) <= len(stringStorage) {
			origBytes = stringStorage[rec.Offset : rec.Offset+rec.Length]
		}
		// Update only Family Name records.
		if rec.NameID == 1 || rec.NameID == 16 {
			encoded, err := encodeName(rec, newFamily)
			if err != nil {
				return nil, err
			}
			newRec.Length = uint16(len(encoded))
			origBytes = encoded
		}
		newRec.Offset = uint16(newBlock.Len())
		newBlock.Write(origBytes)
		newRecords[i] = newRec
	}

	var buf bytes.Buffer
	newHeader := NameTableHeader{
		Format:       header.Format,
		Count:        header.Count,
		StringOffset: uint16(6 + 12*count),
	}

	if err := binary.Write(&buf, binary.BigEndian, newHeader); err != nil {
		return nil, err
	}
	if err := writeStructs(&buf, newRecords); err != nil {
		return nil, err
	}

	buf.Write(newBlock.Bytes())
	return buf.Bytes(), nil
}

func encodeName(rec NameRecord, name string) ([]byte, error) {
	switch rec.PlatformID {
	case 3, 0:
		u16 := utf16.Encode([]rune(name))
		b := make([]byte, len(u16)*2)
		for i, v := range u16 {
			binary.BigEndian.PutUint16(b[i*2:], uint16(v))
		}
		return b, nil
	case 1:
		return []byte(name), nil
	default:
		return nil, errors.New("unsupported platform encoding")
	}
}
