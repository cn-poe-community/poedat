package dat

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// https://github.com/SnosMe/poe-dat-viewer/blob/master/lib/src/dat/dat-file.ts

const rowCountSize = 4 // 4 bytes for row count

var vdataMagic = []byte{0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb}
var minFileSize = rowCountSize + len(vdataMagic)

// |--------|-------------|------------------------|............
// row count|  DataFixed  |  boundary(vdataMagic)  |  DataVariable

type DatFile struct {
	MemSize      int
	RowCount     int
	RowLength    int
	DataFixed    []byte
	DataVariable []byte
	FieldSize    FieldSize
}

func ReadDatFile(content []byte) (*DatFile, error) {
	if len(content) < minFileSize {
		return nil, errors.New("invalid file size")
	}

	rowCount := binary.LittleEndian.Uint32(content[:rowCountSize])
	boundaryIdx := bytes.Index(content, vdataMagic)
	if boundaryIdx == -1 {
		return nil, errors.New("invalid file: section with variable data not found")
	}

	rowLen := 0
	if rowCount > 0 {
		rowLen = (boundaryIdx - rowCountSize) / int(rowCount)
	}

	dataFiexed := content[rowCountSize:boundaryIdx]
	dataVariable := content[boundaryIdx:]

	return &DatFile{
		MemSize:      8,
		RowCount:     int(rowCount),
		RowLength:    rowLen,
		DataFixed:    dataFiexed,
		DataVariable: dataVariable,
		FieldSize:    DefaultFieldSize(),
	}, nil
}
