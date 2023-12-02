package dat

import (
	"bytes"
	"encoding/binary"
	"log"

	binaryutil "github.com/cn-poe-community/poedat/utils/binary"
	"golang.org/x/text/encoding/unicode"
)

//https://github.com/SnosMe/poedat-viewer/blob/master/lib/src/dat/reader.ts

var StringTerminator = []byte{0x00, 0x00, 0x00, 0x00}

const Mem32Null = 0xfefefefe

type FieldSizes struct {
	Bool       int
	String     int
	Key        int
	KeyForeign int
	Array      int
}

func DefaultFieldSize() FieldSizes {
	return FieldSizes{
		Bool:       1,
		String:     8,
		Key:        8,
		KeyForeign: 16,
		Array:      16,
	}
}

// boolean,string,int64,uint64,float64,array,key,keyforeign
type FieldValue interface{}

func ReadUint32(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
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

func ReadKey(b []byte) *uint32 {
	n := binary.LittleEndian.Uint32(b)
	if n == Mem32Null {
		return nil
	} else {
		return &n
	}
}

func ReadKeyForeign(b []byte) *uint32 {
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

func ReadOne(b []byte, h *Header, datFile *DatFile) FieldValue {
	dataVariable := datFile.DataVariable

	t := h.Type
	if t.Boolean {
		return b[0] > 0
	} else if t.String {
		return ReadString(b, dataVariable)
	} else if t.Integer != nil {
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
		return datFile.FieldSizes.Bool
	} else if t.String {
		return datFile.FieldSizes.String
	} else if t.Integer != nil {
		return t.Integer.Size
	} else if t.Decimal != nil {
		return t.Decimal.Size
	} else if t.Key != nil {
		if t.Key.Foreign {
			return datFile.FieldSizes.KeyForeign
		} else {
			return datFile.FieldSizes.Key
		}
	}
	return 0
}

func ReadArray(b []byte, header *Header, dataFile *DatFile) []FieldValue {
	dataVariable := dataFile.DataVariable

	arrayLen := int(ReadUint32(b))
	if arrayLen == 0 {
		return nil
	}

	varOffset := int(ReadUint32(b[dataFile.MemSize:]))
	values := []FieldValue{}
	elementSize := ElementSize(header, dataFile)
	for i := 0; i < arrayLen; i++ {
		var value FieldValue = nil
		value = ReadOne(dataVariable[varOffset+i*elementSize:], header, dataFile)
		values = append(values, value)
	}

	return values
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

func ReadField(b []byte, header *Header, datFile *DatFile) FieldValue {
	t := header.Type

	if t.Array {
		return ReadArray(b, header, datFile)
	} else {
		return ReadOne(b, header, datFile)
	}
}

func ReadRow(rowIndex int, headers []*Header, datFile *DatFile) []FieldValue {
	rowBasicOffset := rowIndex * datFile.RowLength

	values := []FieldValue{}
	for _, header := range headers {
		offset := rowBasicOffset + header.Offset
		values = append(values, ReadField(datFile.DataFixed[offset:], header, datFile))
	}
	return values
}
