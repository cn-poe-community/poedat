package dat

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const rowCountSize = 4

var vdataMagic = []byte{0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb}
var minFileSize = rowCountSize + len(vdataMagic)

// |--------|--------------|--------|--------------------------
// row count|              |boundary
//          |  DataFiexed  |            DataVariable

type DatFile struct {
	MemSize      int
	RowCount     int
	RowLength    int
	DataFixed    []byte
	DataVariable []byte
	FieldSizes   FieldSizes
}

func ReadDatFile(content []byte) (*DatFile, error) {
	if len(content) < minFileSize {
		return nil, errors.New("invalid file size")
	}

	rowCount := binary.LittleEndian.Uint32(content[:rowCountSize])
	boundary := bytes.Index(content, vdataMagic)
	if boundary == -1 {
		return nil, errors.New("invalid file: section with variable data not found")
	}

	rowLen := 0
	if rowCount > 0 {
		rowLen = (boundary - rowCountSize) / int(rowCount)
	}

	dataFiexed := content[rowCountSize:boundary]
	dataVariable := content[boundary:]

	return &DatFile{
		MemSize:      8,
		RowCount:     int(rowCount),
		RowLength:    rowLen,
		DataFixed:    dataFiexed,
		DataVariable: dataVariable,
		FieldSizes:   DefaultFieldSize(),
	}, nil
}
