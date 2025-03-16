package internal

import (
    "bytes"
    "encoding/binary"
    "errors"
    "fmt"
    "math/bits"
    "os"
    "path/filepath"
    "strings"
)

func ProcessFontFile(filePath, newFamily string) error {
    origData, err := os.ReadFile(filePath)
    if err != nil {
        return err
    }
    r := bytes.NewReader(origData)
    offTable, err := readStruct[OffsetTable](r)
    if err != nil {
        return err
    }

    numTables := int(offTable.NumTables)
    origRecords := make([]TableRecord, numTables)
    for i := 0; i < numTables; i++ {
        rec, err := readStruct[TableRecord](r)
        if err != nil {
            return err
        }
        origRecords[i] = rec
    }

    tables := make(map[string]*Table)
    var order []string
    for _, rec := range origRecords {
        tag := string(rec.Tag[:])
        start, end := int(rec.Offset), int(rec.Offset)+int(rec.Length)
        if end > len(origData) {
            return fmt.Errorf("table %s out of bounds", tag)
        }
        // Create a copy of the table data.
        tables[tag] = &Table{Tag: tag, Data: append([]byte{}, origData[start:end]...)}
        order = append(order, tag)
    }

    // Update only Font Family identifiers (NameID 1 and 16).
    nameTable, ok := tables["name"]
    if !ok {
        return errors.New("name table not found")
    }
    newNameTable, err := rebuildNameTable(nameTable.Data, newFamily)
    if err != nil {
        return err
    }
    tables["name"].Data = newNameTable

    newData, newRecords := rebuildTTF(offTable, tables, order)
    newData, err = updateHeadChecksumAdjustment(newData, newRecords)
    if err != nil {
        return err
    }

    dir, file := filepath.Split(filePath)
    ext := filepath.Ext(file)
    base := strings.TrimSuffix(file, ext)
    newFile := filepath.Join(dir, base+"_renamed"+ext)
    if err := os.WriteFile(newFile, newData, 0644); err != nil {
        return err
    }
    fmt.Printf("Written: %s\n", newFile)

    return nil
}

func rebuildTTF(offTable OffsetTable, tables map[string]*Table, order []string) ([]byte, []TableRecord) {
    numTables := len(order)
    headerSize := 12 + 16*numTables
    newRecords := make([]TableRecord, numTables)
    var tableDatas [][]byte
    offset := headerSize

    for i, tag := range order {
        data := tables[tag].Data
        padded := padTo4(data)
        var rec TableRecord
        copy(rec.Tag[:], []byte(tag))
        rec.Offset = uint32(offset)
        rec.Length = uint32(len(data))
        rec.CheckSum = computeChecksum(padded)
        newRecords[i] = rec
        tableDatas = append(tableDatas, padded)
        offset += len(padded)
    }

    maxPow2 := 1
    for maxPow2*2 <= numTables {
        maxPow2 *= 2
    }
    offTable.SearchRange = uint16(maxPow2 * 16)
    offTable.EntrySelector = uint16(bits.Len(uint(maxPow2)) - 1)
    offTable.RangeShift = uint16(numTables*16) - offTable.SearchRange

    var buf bytes.Buffer
    if err := binary.Write(&buf, binary.BigEndian, offTable); err != nil {
        return nil, nil
    }
    if err := writeStructs(&buf, newRecords); err != nil {
        return nil, nil
    }
    for _, d := range tableDatas {
        buf.Write(d)
    }

    return buf.Bytes(), newRecords
}

func padTo4(data []byte) []byte {
    pad := (4 - (len(data) % 4)) % 4
    return append(data, make([]byte, pad)...)
}

func computeChecksum(data []byte) uint32 {
    var sum uint32
    data = padTo4(data)
    for i := 0; i < len(data); i += 4 {
        sum += binary.BigEndian.Uint32(data[i : i+4])
    }

    return sum
}

func updateHeadChecksumAdjustment(ttfData []byte, records []TableRecord) ([]byte, error) {
    var headRec *TableRecord
    for i := range records {
        if tag := string(records[i].Tag[:]); tag == "head" {
            headRec = &records[i]
            break
        }
    }

    if headRec == nil {
        return nil, errors.New("head table not found")
    }
    headStart := int(headRec.Offset)
    if headStart+12 > len(ttfData) {
        return nil, errors.New("head table too short")
    }

    copy(ttfData[headStart+8:headStart+12], []byte{0, 0, 0, 0})
    sum := computeChecksum(ttfData)
    const magic uint32 = 0xB1B0AFBA
    binary.BigEndian.PutUint32(ttfData[headStart+8:headStart+12], magic-sum)

    return ttfData, nil
}
