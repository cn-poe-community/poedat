package dat

import "log"

// https://github.com/SnosMe/poe-dat-viewer/blob/master/lib/src/dat/header.ts

type Header struct {
	Name   string
	Offset int
	Type   HeaderType
}

type HeaderType struct {
	Array    bool
	Interval bool
	Boolean  bool
	Integer  *IntergerType
	Decimal  *DecimalType
	String   bool
	Key      *KeyType
}

type IntergerType struct {
	Unsigned bool
	Size     int
}

type DecimalType struct {
	Size int
}

type KeyType struct {
	Foreign bool
}

func (h *Header) HeaderLength(f FieldSize) int {
	t := &h.Type

	count := 1
	if t.Interval {
		count = 2
	}

	if t.Array {
		return f.Array
	}
	if t.String {
		return f.String
	}
	if t.Key != nil {
		if t.Key.Foreign {
			return f.KeyForeign
		} else {
			return f.Key
		}
	}
	if t.Integer != nil {
		return t.Integer.Size * count
	}
	if t.Decimal != nil {
		return t.Decimal.Size * count
	}
	if t.Boolean {
		return f.Bool
	}

	log.Fatal("Corrupted header")
	return 0
}
