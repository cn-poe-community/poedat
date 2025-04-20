package dat

import (
	"bytes"
	"encoding/binary"
	"log"

	binaryutil "github.com/cn-poe-community/poedat/utils/binary"
	"golang.org/x/text/encoding/unicode"
)

//https://github.com/SnosMe/poe-dat-viewer/blob/master/lib/src/dat/reader.ts

const Mem32Null = 0xfefefefe

var StringTerminator = []byte{0x00, 0x00, 0x00, 0x00}

type FieldSize struct {
	Bool       int
	String     int
	Key        int
	KeyForeign int
	Array      int
}

func DefaultFieldSize() FieldSize {
	return FieldSize{
		Bool:       1,
		String:     8,
		Key:        8,
		KeyForeign: 16,
		Array:      16,
	}
}

type Key *uint32
type KeyForeign *uint32

// bool | string | int16 | int32 | uint16 | uint32 | Key | KeyForeign |
// [](bool | string | int16 | int32 | uint16 | uint32 | Key | KeyForeign) |
// [2](bool | string | int16 | int32 | uint16 | uint32 | Key | KeyForeign)
type FieldValue any

func ReadUint16(b []byte) uint16 {
	return binary.LittleEndian.Uint16(b)
}

func ReadUint32(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

func ReadInt16(b []byte) int16 {
	return int16(binary.LittleEndian.Uint16(b))
}

func ReadInt32(b []byte) int32 {
	return int32(binary.LittleEndian.Uint32(b))
}

func ReadFloat32(b []byte) float32 {
	n, err := binaryutil.LittleEndian.Float32(b)
	if err != nil {
		log.Fatalf("read float32 failed: %v", err)
	}
	return n
}

func ReadKey(b []byte) Key {
	n := binary.LittleEndian.Uint32(b)
	if n == Mem32Null {
		return nil
	} else {
		return &n
	}
}

func ReadKeyForeign(b []byte) KeyForeign {
	n := binary.LittleEndian.Uint32(b)
	if n == Mem32Null {
		return nil
	} else {
		return &n
	}
}

func ReadString(b []byte, dataVariable []byte) string {
	// string stores on dataVariable
	varOffset := int(binary.LittleEndian.Uint32(b))
	end := indexStringTerminator(dataVariable[varOffset:])
	if end == -1 {
		log.Fatalf("read string failed: no matched string terminator ")
	}
	utf8bytes, err := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder().Bytes(dataVariable[varOffset : end+varOffset])
	if err != nil {
		log.Fatalf("read string failed: %v", err)
	}
	return string(utf8bytes)
}

func indexStringTerminator(data []byte) int {
	begin := 0
	for {
		i := bytes.Index(data[begin:], StringTerminator)
		if i == -1 {
			return i
		}
		if (begin+i)%2 == 0 {
			return begin + i
		}
		begin += i + 1
	}
}

func ReadOne(b []byte, h *Header, datFile *DatFile) FieldValue {
	dataVariable := datFile.DataVariable

	t := h.Type
	if t.Boolean {
		return b[0] > 0
	} else if t.String {
		return ReadString(b, dataVariable)
	} else if t.Integer != nil {
		if t.Integer.Unsigned && t.Integer.Size == 2 {
			return ReadUint16(b)
		}
		if t.Integer.Unsigned && t.Integer.Size == 4 {
			return ReadUint32(b)
		}
		if !t.Integer.Unsigned && t.Integer.Size == 2 {
			return ReadInt16(b)
		}
		if !t.Integer.Unsigned && t.Integer.Size == 4 {
			return ReadInt32(b)
		}
	} else if t.Decimal != nil {
		if t.Decimal.Size == 4 {
			return ReadFloat32(b)
		}
	} else if t.Key != nil {
		if t.Key.Foreign {
			return ReadKeyForeign(b)
		} else {
			return ReadKey(b)
		}
	}

	log.Fatalf("unhandled header type %v", t)
	return nil
}

func ElementSize(h *Header, datFile *DatFile) int {
	t := h.Type
	if t.Boolean {
		return datFile.FieldSize.Bool
	} else if t.String {
		return datFile.FieldSize.String
	} else if t.Integer != nil {
		return t.Integer.Size
	} else if t.Decimal != nil {
		return t.Decimal.Size
	} else if t.Key != nil {
		if t.Key.Foreign {
			return datFile.FieldSize.KeyForeign
		} else {
			return datFile.FieldSize.Key
		}
	}
	return 0
}

func ReadArray(b []byte, header *Header, datFile *DatFile) []FieldValue {
	dataVariable := datFile.DataVariable

	arrayLen := int(ReadUint32(b))
	if arrayLen == 0 {
		return nil
	}

	varOffset := int(ReadUint32(b[datFile.MemSize:]))
	values := []FieldValue{}
	elementSize := ElementSize(header, datFile)
	for i := 0; i < arrayLen; i++ {
		var value FieldValue = nil
		value = ReadOne(dataVariable[varOffset+i*elementSize:], header, datFile)
		values = append(values, value)
	}

	return values
}

func ReadInterval(b []byte, header *Header, datFile *DatFile) [2]FieldValue {
	elementSize := ElementSize(header, datFile)
	return [2]FieldValue{
		ReadOne(b, header, datFile),
		ReadOne(b[elementSize:], header, datFile),
	}
}

func ReadField(b []byte, header *Header, datFile *DatFile) FieldValue {
	t := header.Type

	if t.Array {
		return ReadArray(b, header, datFile)
	} else if t.Interval {
		return ReadInterval(b, header, datFile)
	} else {
		return ReadOne(b, header, datFile)
	}
}

func ReadRow(rowIdx int, headers []*Header, datFile *DatFile) []FieldValue {
	rowBasicOffset := rowIdx * datFile.RowLength

	values := []FieldValue{}
	for _, header := range headers {
		offset := rowBasicOffset + header.Offset
		values = append(values, ReadField(datFile.DataFixed[offset:], header, datFile))
	}
	return values
}
