package schema

// https://github.com/poe-tool-dev/dat-schema/blob/main/src/types.ts

type ColumnType string

const TypeBool ColumnType = "bool"
const TypeString ColumnType = "string"
const TypeI16 ColumnType = "i16"
const TypeI32 ColumnType = "i32"
const TypeU16 ColumnType = "u16"
const TypeU32 ColumnType = "u32"
const TypeF32 ColumnType = "f32"
const TypeArray ColumnType = "array"
const TypeRow ColumnType = "row"
const TypeForeignRow ColumnType = "foreignRow"
const TypeEnumRow ColumnType = "enumRow"

type Ref struct {
	Table  string  `json:"table"`
	Column *string `json:"column"`
}

type ValidFor int

const ValidForPoe1 ValidFor = 1
const ValidForPoe2 ValidFor = 2
const ValidForCommon ValidFor = 3

type TableColumn struct {
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Array       bool       `json:"array"`
	Type        ColumnType `json:"type"`
	Unique      bool       `json:"unique"`
	Localized   bool       `json:"localized"`
	Util        *string    `json:"util"`
	References  *Ref       `json:"references"`
}

type SchemaTable struct {
	ValidFor ValidFor       `json:"validFor"`
	Name     string         `json:"name"`
	Columns  []*TableColumn `json:"columns"`
	Tags     []string       `json:"tags"`
}

type SchemaEnumeration struct {
	Name        string    `json:"name"`
	Indexing    int       `json:"indexing"`
	Enumerators []*string `json:"enumerators"`
}

type SchemaMetadata struct {
	Version   int `json:"version"`
	CreatedAt int `json:"createdAt"`
}

type SchemaFile struct {
	Version      int                  `json:"version"`
	CreatedAt    int                  `json:"createdAt"`
	Tables       []*SchemaTable       `json:"tables"`
	Enumerations []*SchemaEnumeration `json:"enumerations"`
}
