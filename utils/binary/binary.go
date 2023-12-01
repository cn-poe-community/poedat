package binaryutil

import (
	"bytes"
	"encoding/binary"
)

var LittleEndian littleEndian

type littleEndian struct{}

func (littleEndian) Float32(b []byte) (float32, error) {
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	var n float32
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &n)
	return n, err
}
