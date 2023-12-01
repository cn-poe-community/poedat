package dat

import "log"

type Header struct {
	Name   string
	Offset int
	Type   HeaderType
}

type HeaderType struct {
	Array   bool
	Boolean bool
	Integer *IntergerType
	Decimal *DecimalType
	String  bool
	Key     *KeyType
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

func (h *Header) HeaderLength(sizes FieldSizes) int {
	t := &h.Type
	if t.Array {
		return sizes.Array
	}
	if t.String {
		return sizes.String
	}
	if t.Key != nil {
		if t.Key.Foreign {
			return sizes.KeyForeign
		} else {
			return sizes.Key
		}
	}
	if t.Integer != nil {
		return t.Integer.Size
	}
	if t.Decimal != nil {
		return t.Decimal.Size
	}
	if t.Boolean {
		return sizes.Bool
	}

	log.Fatal("Corrupted header")
	return 0
}
